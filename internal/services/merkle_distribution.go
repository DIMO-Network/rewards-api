package services

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/merkletree"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
)

// TreeUploader uploads a serialized Merkle tree file to object storage under
// the given key.
type TreeUploader interface {
	Upload(ctx context.Context, key string, body []byte) error
}

// S3TreeUploader implements TreeUploader on top of an S3 bucket.
type S3TreeUploader struct {
	Client *s3.Client
	Bucket string
}

func (u *S3TreeUploader) Upload(ctx context.Context, key string, body []byte) error {
	_, err := u.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String("application/json"),
	})
	return err
}

// MerkleDistributionService builds the weekly Merkle tree of reward claims,
// uploads the tree file to object storage, and submits a setRoot meta-transaction
// to the MerkleDistributor contract via the same Kafka path used for push
// transfers.
type MerkleDistributionService struct {
	TransferService    *TransferService
	Uploader           TreeUploader
	DistributorAddress common.Address
	PoolID             int
	BaseURI            string
	Logger             *zerolog.Logger
}

// NewMerkleDistributionService validates the Merkle-related settings and
// constructs a MerkleDistributionService.
func NewMerkleDistributionService(
	settings *config.Settings,
	transferService *TransferService,
	uploader TreeUploader,
	logger *zerolog.Logger,
) (*MerkleDistributionService, error) {
	if !common.IsHexAddress(settings.MerkleDistributorAddress) {
		return nil, fmt.Errorf("invalid MERKLE_DISTRIBUTOR_ADDRESS %q", settings.MerkleDistributorAddress)
	}
	if settings.MerklePoolID < 0 {
		return nil, fmt.Errorf("MERKLE_POOL_ID must be non-negative, got %d", settings.MerklePoolID)
	}
	if settings.MerkleTreeBaseURI == "" {
		return nil, errors.New("MERKLE_TREE_BASE_URI must be set")
	}

	return &MerkleDistributionService{
		TransferService:    transferService,
		Uploader:           uploader,
		DistributorAddress: common.HexToAddress(settings.MerkleDistributorAddress),
		PoolID:             settings.MerklePoolID,
		BaseURI:            settings.MerkleTreeBaseURI,
		Logger:             logger,
	}, nil
}

// DistributeWeek builds the Merkle tree for the given issuance week, uploads
// the tree file, and submits a setRoot meta-transaction. It is idempotent: if
// the root has already been set successfully, or a meta-transaction is still
// pending, it does nothing.
func (m *MerkleDistributionService) DistributeWeek(ctx context.Context, week int) error {
	logger := m.Logger.With().Int("issuanceWeek", week).Int("poolId", m.PoolID).Logger()

	existing, err := models.FindMerkleRoot(ctx, m.TransferService.db.DBS().Reader, week)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to look up existing merkle root: %w", err)
	}
	if existing != nil {
		if existing.SetSuccessful {
			logger.Info().Msg("Merkle root already set successfully for this week. Nothing to do.")
			return nil
		}
		if existing.MetaTransactionRequestID.Valid {
			logger.Info().Str("requestId", existing.MetaTransactionRequestID.String).Msg("Merkle root meta-transaction still pending for this week. Nothing to do.")
			return nil
		}
	}

	rewards, err := models.Rewards(
		qm.Expr(
			qm.Or2(models.RewardWhere.SyntheticDeviceTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
			qm.Or2(models.RewardWhere.AftermarketDeviceTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
			qm.Or2(models.RewardWhere.StreakTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
		),
		models.RewardWhere.IssuanceWeekID.EQ(week),
		// Temporary blacklist, see PLA-765.
		qm.LeftOuterJoin("rewards_api."+models.TableNames.Blacklist+" ON "+models.BlacklistTableColumns.UserEthereumAddress+" = "+models.RewardTableColumns.UserEthereumAddress),
		qm.Where(models.BlacklistTableColumns.UserEthereumAddress+" IS NULL"),
	).All(ctx, m.TransferService.db.DBS().Reader)
	if err != nil {
		return fmt.Errorf("failed to load rewards for week %d: %w", week, err)
	}

	leaves, totalAllocation := aggregateRewardLeaves(rewards, &logger)
	if len(leaves) == 0 {
		logger.Warn().Msg("No positive reward balances for this week. Not setting a Merkle root.")
		return nil
	}

	tree, err := merkletree.New(m.DistributorAddress, big.NewInt(int64(m.PoolID)), big.NewInt(int64(week)), leaves)
	if err != nil {
		return fmt.Errorf("failed to build merkle tree: %w", err)
	}

	treeJSON, err := tree.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize merkle tree: %w", err)
	}

	key := TreeFileKey(m.PoolID, week)
	if err := m.Uploader.Upload(ctx, key, treeJSON); err != nil {
		return fmt.Errorf("failed to upload merkle tree file %s: %w", key, err)
	}
	proofsURI := strings.TrimSuffix(m.BaseURI, "/") + "/" + key

	root := tree.Root()
	logger.Info().
		Str("root", root.Hex()).
		Str("totalAllocation", totalAllocation.String()).
		Int("leaves", len(leaves)).
		Str("proofsURI", proofsURI).
		Msg("Built and uploaded merkle tree.")

	calldata, err := packSetRoot(m.PoolID, week, root, totalAllocation, proofsURI)
	if err != nil {
		return fmt.Errorf("failed to pack setRoot call: %w", err)
	}

	tx, err := m.TransferService.db.DBS().Writer.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint

	reqID := ksuid.New().String()
	metaTxRequest := &models.MetaTransactionRequest{
		ID:     reqID,
		Status: models.MetaTransactionRequestStatusUnsubmitted,
	}
	if err := metaTxRequest.Insert(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert meta-transaction request: %w", err)
	}

	rootRow := &models.MerkleRoot{
		IssuanceWeekID:           week,
		PoolID:                   m.PoolID,
		Root:                     root.Bytes(),
		TotalAllocation:          types.NewDecimal(new(decimal.Big).SetBigMantScale(totalAllocation, 0)),
		ProofsURI:                proofsURI,
		MetaTransactionRequestID: null.StringFrom(reqID),
	}
	if err := rootRow.Upsert(ctx, tx, true, []string{models.MerkleRootColumns.IssuanceWeekID}, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("failed to upsert merkle root row: %w", err)
	}

	// Send the Kafka request before committing: a send failure rolls back both
	// rows, so no orphaned pending request is left behind. The remaining crash
	// window (send succeeds, commit fails) can produce a Kafka message with no
	// matching rows, which is tolerated because the meta-transaction processor
	// is idempotent by request ID and the status listener ignores unknown IDs.
	if err := m.TransferService.sendRequest(reqID, m.DistributorAddress, calldata); err != nil {
		return fmt.Errorf("failed to send setRoot meta-transaction request: %w", err)
	}

	return tx.Commit()
}

// TreeFileKey returns the object storage key for the tree file of the given
// pool and week.
func TreeFileKey(poolID, week int) string {
	return fmt.Sprintf("pool-%d/week-%d.json", poolID, week)
}

// aggregateRewardLeaves sums the aftermarket, synthetic, and streak token
// amounts of the given reward rows per rewards receiver address, returning one
// leaf per address with a positive total, sorted by address, together with the
// sum of all leaf amounts.
func aggregateRewardLeaves(rewards models.RewardSlice, logger *zerolog.Logger) ([]merkletree.Leaf, *big.Int) {
	amountByAccount := make(map[common.Address]*big.Int)

	for _, row := range rewards {
		if !row.RewardsReceiverEthereumAddress.Valid || !common.IsHexAddress(row.RewardsReceiverEthereumAddress.String) {
			logger.Warn().Int("vehicleId", row.UserDeviceTokenID).Msg("Reward row has no valid rewards receiver address. Skipping.")
			continue
		}
		account := common.HexToAddress(row.RewardsReceiverEthereumAddress.String)

		rowTotal := new(big.Int).Add(
			utils.NullDecimalToIntDefaultZero(row.AftermarketDeviceTokens),
			new(big.Int).Add(
				utils.NullDecimalToIntDefaultZero(row.SyntheticDeviceTokens),
				utils.NullDecimalToIntDefaultZero(row.StreakTokens),
			),
		)
		if rowTotal.Sign() <= 0 {
			continue
		}

		if existing, ok := amountByAccount[account]; ok {
			existing.Add(existing, rowTotal)
		} else {
			amountByAccount[account] = rowTotal
		}
	}

	leaves := make([]merkletree.Leaf, 0, len(amountByAccount))
	totalAllocation := big.NewInt(0)
	for account, amount := range amountByAccount {
		leaves = append(leaves, merkletree.Leaf{Account: account, Amount: amount})
		totalAllocation.Add(totalAllocation, amount)
	}
	sort.Slice(leaves, func(i, j int) bool {
		return bytes.Compare(leaves[i].Account[:], leaves[j].Account[:]) < 0
	})

	return leaves, totalAllocation
}

// packSetRoot ABI-encodes a call to MerkleDistributor.setRoot.
func packSetRoot(poolID, week int, root common.Hash, allocation *big.Int, proofsURI string) ([]byte, error) {
	distributorABI, err := contracts.MerkleDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return distributorABI.Pack("setRoot", big.NewInt(int64(poolID)), big.NewInt(int64(week)), [32]byte(root), allocation, proofsURI)
}

// ValidateMerkleCutover guards against double payment: if the Merkle path is
// enabled, the first Merkle week must be strictly greater than every week that
// has already been paid out via push transfers.
func ValidateMerkleCutover(ctx context.Context, dbs db.Store, firstMerkleWeek int) error {
	if firstMerkleWeek <= 0 {
		return nil
	}

	var result struct {
		MaxWeek null.Int `boil:"max_week"`
	}
	err := models.NewQuery(
		qm.Select("max("+models.RewardColumns.IssuanceWeekID+") AS max_week"),
		qm.From("rewards_api."+models.TableNames.Rewards),
		qm.Where(models.RewardColumns.TransferSuccessful+" = true"),
	).Bind(ctx, dbs.DBS().Reader, &result)
	if err != nil {
		return fmt.Errorf("failed to determine last successfully transferred week: %w", err)
	}

	if result.MaxWeek.Valid && firstMerkleWeek <= result.MaxWeek.Int {
		return fmt.Errorf("FIRST_MERKLE_WEEK %d is not after the last successfully push-transferred week %d; refusing to run to avoid double payment", firstMerkleWeek, result.MaxWeek.Int)
	}

	return nil
}
