package services

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/DIMO-Network/cloudevent"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/merkletree"
	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func nullDec(i int64) types.NullDecimal {
	return types.NewNullDecimal(decimal.New(i, 0))
}

func TestAggregateRewardLeaves(t *testing.T) {
	logger := zerolog.Nop()

	addr1 := mkAddr(1)
	addr2 := mkAddr(2)

	rewards := models.RewardSlice{
		// addr1 has two vehicles; amounts must be summed into one leaf.
		{
			UserDeviceTokenID:              1,
			RewardsReceiverEthereumAddress: null.StringFrom(addr1.Hex()),
			AftermarketDeviceTokens:        nullDec(100),
			StreakTokens:                   nullDec(50),
		},
		{
			UserDeviceTokenID:              2,
			RewardsReceiverEthereumAddress: null.StringFrom(addr1.Hex()),
			SyntheticDeviceTokens:          nullDec(25),
		},
		{
			UserDeviceTokenID:              3,
			RewardsReceiverEthereumAddress: null.StringFrom(addr2.Hex()),
			SyntheticDeviceTokens:          nullDec(40),
		},
		// Zero amounts must not produce a leaf.
		{
			UserDeviceTokenID:              4,
			RewardsReceiverEthereumAddress: null.StringFrom(mkAddr(3).Hex()),
		},
		// No receiver address must not produce a leaf.
		{
			UserDeviceTokenID:       5,
			AftermarketDeviceTokens: nullDec(7),
		},
	}

	leaves, total := aggregateRewardLeaves(rewards, &logger)

	assert.Equal(t, []merkletree.Leaf{
		{Account: addr1, Amount: big.NewInt(175)},
		{Account: addr2, Amount: big.NewInt(40)},
	}, leaves)
	assert.Equal(t, big.NewInt(215), total)
}

type fakeUploader struct {
	uploads map[string][]byte
}

func (f *fakeUploader) Upload(_ context.Context, key string, body []byte) error {
	if f.uploads == nil {
		f.uploads = make(map[string][]byte)
	}
	f.uploads[key] = body
	return nil
}

func merkleTestSettings() *config.Settings {
	return &config.Settings{
		MetaTransactionSendTopic: "topic.transaction.request.send",
		TransferBatchSize:        100,
		MerkleDistributorAddress: "0x00000000000000000000000000000000000000aA",
		MerklePoolID:             0,
		MerkleTreeBaseURI:        "https://merkle.dimo.zone/",
		MerkleTreeS3Bucket:       "merkle-trees-test",
		FirstMerkleWeek:          100,
	}
}

func TestMerkleDistributeWeek(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout)

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	settings := merkleTestSettings()

	kafkaConfig := mocks.NewTestConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	producer := mocks.NewSyncProducer(t, kafkaConfig)

	transferService := NewTokenTransferService(settings, producer, conn)

	uploader := &fakeUploader{}
	svc, err := NewMerkleDistributionService(settings, transferService, uploader, &logger)
	require.NoError(t, err)

	var sent []cloudevent.CloudEvent[transferData]
	checker := func(b []byte) error {
		var o cloudevent.CloudEvent[transferData]
		require.NoError(t, json.Unmarshal(b, &o))
		sent = append(sent, o)
		return nil
	}

	const week = 100

	wk := models.IssuanceWeek{ID: week, JobStatus: models.IssuanceWeeksJobStatusFinished}
	require.NoError(t, wk.Insert(ctx, conn.DBS().Writer, boil.Infer()))

	addr1 := mkAddr(1)
	addr2 := mkAddr(2)
	blacklisted := mkAddr(3)

	rewards := []models.Reward{
		{
			IssuanceWeekID:                 week,
			UserDeviceTokenID:              1,
			UserEthereumAddress:            null.StringFrom(addr1.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(addr1.Hex()),
			AftermarketDeviceTokens:        nullDec(100),
			StreakTokens:                   nullDec(50),
		},
		{
			IssuanceWeekID:                 week,
			UserDeviceTokenID:              2,
			UserEthereumAddress:            null.StringFrom(addr1.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(addr1.Hex()),
			SyntheticDeviceTokens:          nullDec(25),
		},
		{
			IssuanceWeekID:                 week,
			UserDeviceTokenID:              3,
			UserEthereumAddress:            null.StringFrom(addr2.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(addr2.Hex()),
			SyntheticDeviceTokens:          nullDec(40),
		},
		// Blacklisted user must be excluded.
		{
			IssuanceWeekID:                 week,
			UserDeviceTokenID:              4,
			UserEthereumAddress:            null.StringFrom(blacklisted.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(blacklisted.Hex()),
			SyntheticDeviceTokens:          nullDec(999),
		},
		// Zero amounts must be excluded.
		{
			IssuanceWeekID:                 week,
			UserDeviceTokenID:              5,
			UserEthereumAddress:            null.StringFrom(mkAddr(4).Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(mkAddr(4).Hex()),
		},
	}
	for i := range rewards {
		require.NoError(t, rewards[i].Insert(ctx, conn.DBS().Writer, boil.Infer()))
	}

	bl := models.Blacklist{UserEthereumAddress: blacklisted.Hex(), Note: "test"}
	require.NoError(t, bl.Insert(ctx, conn.DBS().Writer, boil.Infer()))

	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)
	require.NoError(t, svc.DistributeWeek(ctx, week))

	// Independently built tree must match.
	expectedTree, err := merkletree.New(
		common.HexToAddress(settings.MerkleDistributorAddress),
		big.NewInt(0),
		big.NewInt(week),
		[]merkletree.Leaf{
			{Account: addr1, Amount: big.NewInt(175)},
			{Account: addr2, Amount: big.NewInt(40)},
		},
	)
	require.NoError(t, err)

	rootRow, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, week)
	require.NoError(t, err)

	assert.Equal(t, expectedTree.Root().Bytes(), rootRow.Root)
	assert.Equal(t, 0, rootRow.PoolID)
	assert.Equal(t, "https://merkle.dimo.zone/pool-0/week-100.json", rootRow.ProofsURI)
	assert.False(t, rootRow.SetSuccessful)
	require.True(t, rootRow.MetaTransactionRequestID.Valid)
	assert.Equal(t, "215", rootRow.TotalAllocation.String())

	// The uploaded tree file must parse, verify, and contain both leaves.
	body, ok := uploader.uploads["pool-0/week-100.json"]
	require.True(t, ok, "expected tree file upload")
	treeFile, err := merkletree.UnmarshalTreeFile(body)
	require.NoError(t, err)
	require.NoError(t, treeFile.VerifyRoot())
	assert.Len(t, treeFile.Leaves, 2)
	assert.Equal(t, expectedTree.Root(), treeFile.Root)

	// The Kafka message must target the distributor with a setRoot call
	// carrying the expected arguments.
	require.Len(t, sent, 1)
	assert.Equal(t, rootRow.MetaTransactionRequestID.String, sent[0].Data.ID)
	assert.Equal(t, common.HexToAddress(settings.MerkleDistributorAddress), sent[0].Data.To)

	distributorABI, err := contracts.MerkleDistributorMetaData.GetAbi()
	require.NoError(t, err)
	calldata := []byte(sent[0].Data.Data)
	method, err := distributorABI.MethodById(calldata[:4])
	require.NoError(t, err)
	assert.Equal(t, "setRoot", method.Name)

	args, err := method.Inputs.Unpack(calldata[4:])
	require.NoError(t, err)
	require.Len(t, args, 5)
	assert.Equal(t, "0", args[0].(*big.Int).String())
	assert.Equal(t, "100", args[1].(*big.Int).String())
	assert.Equal(t, [32]byte(expectedTree.Root()), args[2].([32]byte))
	assert.Equal(t, "215", args[3].(*big.Int).String())
	assert.Equal(t, "https://merkle.dimo.zone/pool-0/week-100.json", args[4].(string))

	t.Run("pending request is not resubmitted", func(t *testing.T) {
		// No new producer expectation: an unexpected send would fail the test.
		require.NoError(t, svc.DistributeWeek(ctx, week))
	})

	t.Run("successful root is not resubmitted", func(t *testing.T) {
		rootRow.SetSuccessful = true
		rootRow.MetaTransactionRequestID = null.String{}
		_, err := rootRow.Update(ctx, conn.DBS().Writer, boil.Infer())
		require.NoError(t, err)

		require.NoError(t, svc.DistributeWeek(ctx, week))
	})

	t.Run("failed root is retried with a new request", func(t *testing.T) {
		rootRow.SetSuccessful = false
		rootRow.MetaTransactionRequestID = null.String{}
		_, err := rootRow.Update(ctx, conn.DBS().Writer, boil.Infer())
		require.NoError(t, err)

		producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)
		require.NoError(t, svc.DistributeWeek(ctx, week))

		require.Len(t, sent, 2)
		retried, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, week)
		require.NoError(t, err)
		require.True(t, retried.MetaTransactionRequestID.Valid)
		assert.Equal(t, sent[1].Data.ID, retried.MetaTransactionRequestID.String)
		assert.Equal(t, expectedTree.Root().Bytes(), retried.Root)
	})

	t.Run("kafka send failure leaves no orphan rows", func(t *testing.T) {
		const failWeek = 102
		wk := models.IssuanceWeek{ID: failWeek, JobStatus: models.IssuanceWeeksJobStatusFinished}
		require.NoError(t, wk.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		rw := models.Reward{
			IssuanceWeekID:                 failWeek,
			UserDeviceTokenID:              1,
			UserEthereumAddress:            null.StringFrom(addr1.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(addr1.Hex()),
			StreakTokens:                   nullDec(10),
		}
		require.NoError(t, rw.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		mtrsBefore, err := models.MetaTransactionRequests().Count(ctx, conn.DBS().Reader)
		require.NoError(t, err)

		producer.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)
		require.Error(t, svc.DistributeWeek(ctx, failWeek))

		// The transaction must roll back: no merkle root row and no orphaned
		// meta-transaction request.
		exists, err := models.MerkleRootExists(ctx, conn.DBS().Reader, failWeek)
		require.NoError(t, err)
		assert.False(t, exists)

		mtrsAfter, err := models.MetaTransactionRequests().Count(ctx, conn.DBS().Reader)
		require.NoError(t, err)
		assert.Equal(t, mtrsBefore, mtrsAfter)
	})

	t.Run("week without rewards sets no root", func(t *testing.T) {
		const emptyWeek = 101
		wk := models.IssuanceWeek{ID: emptyWeek, JobStatus: models.IssuanceWeeksJobStatusFinished}
		require.NoError(t, wk.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		require.NoError(t, svc.DistributeWeek(ctx, emptyWeek))

		exists, err := models.MerkleRootExists(ctx, conn.DBS().Reader, emptyWeek)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	require.NoError(t, producer.Close())
}

func TestValidateMerkleCutover(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout)

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	// Disabled Merkle path always passes.
	require.NoError(t, ValidateMerkleCutover(ctx, conn, 0))

	// No push transfers at all: any positive week passes.
	require.NoError(t, ValidateMerkleCutover(ctx, conn, 1))

	wk := models.IssuanceWeek{ID: 90, JobStatus: models.IssuanceWeeksJobStatusFinished}
	require.NoError(t, wk.Insert(ctx, conn.DBS().Writer, boil.Infer()))

	rw := models.Reward{
		IssuanceWeekID:                 90,
		UserDeviceTokenID:              1,
		RewardsReceiverEthereumAddress: null.StringFrom(mkAddr(1).Hex()),
		StreakTokens:                   nullDec(10),
		TransferSuccessful:             null.BoolFrom(true),
	}
	require.NoError(t, rw.Insert(ctx, conn.DBS().Writer, boil.Infer()))

	assert.Error(t, ValidateMerkleCutover(ctx, conn, 90), "first merkle week equal to a paid push week must be rejected")
	assert.Error(t, ValidateMerkleCutover(ctx, conn, 89), "first merkle week before a paid push week must be rejected")
	assert.NoError(t, ValidateMerkleCutover(ctx, conn, 91), "first merkle week after the last paid push week must be accepted")
}

func TestMerkleRootStatusListener(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout)

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	proc, err := NewStatusProcessor(conn, &logger, &config.Settings{})
	require.NoError(t, err)

	seed := func(t *testing.T, week int, requestID string) *models.MerkleRoot {
		t.Helper()

		wk := models.IssuanceWeek{ID: week, JobStatus: models.IssuanceWeeksJobStatusFinished}
		require.NoError(t, wk.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		mtr := models.MetaTransactionRequest{ID: requestID, Status: models.MetaTransactionRequestStatusSubmitted}
		require.NoError(t, mtr.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		root := models.MerkleRoot{
			IssuanceWeekID:           week,
			PoolID:                   0,
			Root:                     common.HexToHash("0x01").Bytes(),
			TotalAllocation:          types.NewDecimal(decimal.New(215, 0)),
			ProofsURI:                "https://merkle.dimo.zone/pool-0/week-100.json",
			MetaTransactionRequestID: null.StringFrom(requestID),
		}
		require.NoError(t, root.Insert(ctx, conn.DBS().Writer, boil.Infer()))

		return &root
	}

	makeMessage := func(t *testing.T, requestID, status string, successful *bool) *sarama.ConsumerMessage {
		t.Helper()

		event := cloudevent.CloudEvent[ceData]{
			Data: ceData{
				RequestID: requestID,
				Type:      status,
				Transaction: ceTx{
					Hash:       "0xabc",
					Successful: successful,
				},
			},
		}
		value, err := json.Marshal(event)
		require.NoError(t, err)
		return &sarama.ConsumerMessage{Value: value}
	}

	boolPtr := func(b bool) *bool { return &b }

	t.Run("confirmed successful flips set_successful", func(t *testing.T) {
		requestID := ksuid.New().String()
		seed(t, 100, requestID)

		require.NoError(t, proc.processMessage(makeMessage(t, requestID, models.MetaTransactionRequestStatusConfirmed, boolPtr(true))))

		root, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, 100)
		require.NoError(t, err)
		assert.True(t, root.SetSuccessful)
		assert.Equal(t, requestID, root.MetaTransactionRequestID.String)

		mtr, err := models.FindMetaTransactionRequest(ctx, conn.DBS().Reader, requestID)
		require.NoError(t, err)
		assert.Equal(t, models.MetaTransactionRequestStatusConfirmed, mtr.Status)
		assert.True(t, mtr.Successful.Bool)
	})

	t.Run("failed clears request id for retry", func(t *testing.T) {
		requestID := ksuid.New().String()
		seed(t, 101, requestID)

		require.NoError(t, proc.processMessage(makeMessage(t, requestID, models.MetaTransactionRequestStatusFailed, nil)))

		root, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, 101)
		require.NoError(t, err)
		assert.False(t, root.SetSuccessful)
		assert.False(t, root.MetaTransactionRequestID.Valid)

		mtr, err := models.FindMetaTransactionRequest(ctx, conn.DBS().Reader, requestID)
		require.NoError(t, err)
		assert.Equal(t, models.MetaTransactionRequestStatusFailed, mtr.Status)
	})

	t.Run("confirmed with nil successful errors and leaves row for redelivery", func(t *testing.T) {
		requestID := ksuid.New().String()
		seed(t, 103, requestID)

		err := proc.processMessage(makeMessage(t, requestID, models.MetaTransactionRequestStatusConfirmed, nil))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing successful field")

		// Nothing must change: the root stays attached to the request and the
		// request status stays untouched, so a Kafka redelivery can retry.
		root, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, 103)
		require.NoError(t, err)
		assert.False(t, root.SetSuccessful)
		require.True(t, root.MetaTransactionRequestID.Valid)
		assert.Equal(t, requestID, root.MetaTransactionRequestID.String)

		mtr, err := models.FindMetaTransactionRequest(ctx, conn.DBS().Reader, requestID)
		require.NoError(t, err)
		assert.Equal(t, models.MetaTransactionRequestStatusSubmitted, mtr.Status)
		assert.False(t, mtr.Successful.Valid)
	})

	t.Run("confirmed reverted clears request id for retry", func(t *testing.T) {
		requestID := ksuid.New().String()
		seed(t, 102, requestID)

		require.NoError(t, proc.processMessage(makeMessage(t, requestID, models.MetaTransactionRequestStatusConfirmed, boolPtr(false))))

		root, err := models.FindMerkleRoot(ctx, conn.DBS().Reader, 102)
		require.NoError(t, err)
		assert.False(t, root.SetSuccessful)
		assert.False(t, root.MetaTransactionRequestID.Valid)

		mtr, err := models.FindMetaTransactionRequest(ctx, conn.DBS().Reader, requestID)
		require.NoError(t, err)
		assert.Equal(t, models.MetaTransactionRequestStatusConfirmed, mtr.Status)
		assert.False(t, mtr.Successful.Bool)
	})
}
