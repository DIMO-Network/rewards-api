package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	mt "github.com/txaty/go-merkletree"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var CredentialExpiresAt = "ExprTime"
var Values = "Values"

type leaf struct {
	data []byte
}

func (l *leaf) Serialize() ([]byte, error) {
	// since OpenZep will sort bytes before hashing
	// abi encoding and serialization is performed before
	// we pass any data to the merkle tree func
	// thus, we can just return the bytes directly here
	return l.data, nil
}

func abiEncode(types []string, values []interface{}) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("type array  and value array must be same length")
	}

	args := abi.Arguments{}
	for _, argType := range types {
		abiType, err := abi.NewType(argType, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, abi.Argument{Type: abiType})
	}

	return args.Pack(values...)
}

type Attestor struct {
	abi       *abi.ABI
	producer  sarama.SyncProducer
	pdb       db.Store
	config    mt.Config
	schema    [32]byte
	recipient common.Address
	contract  common.Address
	mtxTopic  string
	log       *zerolog.Logger
}

func NewAttestor(producer sarama.SyncProducer, pdb db.Store, settings *config.Settings, logger *zerolog.Logger) (*Attestor, error) {
	abi, err := contracts.EasMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	var byteArray [32]byte
	schemaSlice, err := hexutil.Decode(settings.EASSchema)
	if err != nil {
		return nil, err
	}
	copy(byteArray[:], schemaSlice)

	return &Attestor{
		abi:      abi,
		producer: producer,
		pdb:      pdb,
		schema:   byteArray,
		config: mt.Config{
			SortSiblingPairs: true, // parameter for OpenZeppelin compatibility
			HashFunc: func(data []byte) ([]byte, error) {
				return crypto.Keccak256(data), nil
			},
			Mode:               mt.ModeProofGenAndTreeBuild,
			DisableLeafHashing: true, // also required for OpenZeppelin compatilibity
		},
		contract:  common.HexToAddress(settings.AttestationContract),
		recipient: common.HexToAddress(settings.AttestationRecipient),
		mtxTopic:  settings.MetaTransactionSendTopic,
		log:       logger,
	}, nil
}

func (a *Attestor) GenerateAndSubmitAttestation(data map[string]map[string]interface{}, dataTypes []string) (string, error) {
	reqID := ksuid.New().String()
	tx, err := a.pdb.DBS().Writer.BeginTx(context.Background(), nil)
	if err != nil {
		return "", err
	}

	tree, err := a.generateMerkleTree(data, dataTypes)
	if err != nil {
		return "", err
	}

	var root [32]byte
	copy(root[:], tree.Root)
	attestationData, err := a.prepareAttestationRequest(root)
	if err != nil {
		return "", err
	}

	event := shared.CloudEvent[transferData]{
		ID:          reqID,
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     reqID,
		Time:        time.Now(),
		Type:        "zone.dimo.attestation.request",
		Data: transferData{
			ID:   reqID,
			To:   a.contract, // EAS contract
			Data: attestationData,
		},
	}

	b, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	_, _, err = a.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: a.mtxTopic,
			Key:   sarama.StringEncoder(reqID),
			Value: sarama.ByteEncoder(b),
		},
	)
	if err != nil {
		return "", err
	}

	mtr := models.MetaTransactionRequest{
		ID:     reqID,
		Status: models.MetaTransactionRequestStatusSubmitted,
	}
	if err := mtr.Insert(context.Background(), tx, boil.Infer()); err != nil {
		return "", err
	}

	att := models.Attestation{
		TransactionID: reqID,
		Root:          tree.Root,
	}
	if err := att.Insert(context.Background(), tx, boil.Infer()); err != nil {
		return "", err
	}

	return reqID, tx.Commit()
}

func (a *Attestor) generateMerkleTree(data map[string]map[string]interface{}, dataTypes []string) (*mt.MerkleTree, error) {
	encodedLeaves := [][]byte{}
	for _, device := range data {
		if _, val := device[Values]; val {
			codeBytes, err := abiEncode(dataTypes, device[Values].([]interface{}))
			if err != nil {
				a.log.Err(err).Msg("failed to abi encode attestation data")
			}
			encodedLeaves = append(encodedLeaves, crypto.Keccak256(crypto.Keccak256(codeBytes)))
		}
	}

	// sort values for OpenZeppelin compatibility
	sort.Slice(encodedLeaves, func(i, j int) bool {
		return bytes.Compare(encodedLeaves[i], encodedLeaves[j]) < 0
	})

	blocks := []mt.DataBlock{}
	for _, vc := range encodedLeaves {
		blocks = append(blocks, &leaf{data: vc})
	}

	return mt.New(&a.config, blocks)
}

func (a *Attestor) prepareAttestationRequest(root [32]byte) ([]byte, error) {
	args := abi.Arguments{}
	abiType, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}

	args = append(args, abi.Argument{Type: abiType})
	enocodedData, err := args.Pack(root)
	if err != nil {
		return nil, err
	}

	req := contracts.AttestationRequest{
		Schema: a.schema,
		Data: contracts.AttestationRequestData{
			Recipient: a.recipient,
			Data:      enocodedData,
			Value:     big.NewInt(0),
		},
	}
	return a.abi.Pack("attest", req)
}

func (a *Attestor) VerifyAttestation(root, node []byte, proof [][]byte) (bool, error) {
	dataBlock := &leaf{data: node}
	p := mt.Proof{
		Siblings: proof,
	}
	return mt.Verify(dataBlock, &p, root, &a.config)

}

func AbiEncode(types []string, values []interface{}) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("type array  and value array must be same length")
	}

	args := abi.Arguments{}
	for _, argType := range types {
		abiType, err := abi.NewType(argType, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, abi.Argument{Type: abiType})
	}

	return args.Pack(values...)
}
