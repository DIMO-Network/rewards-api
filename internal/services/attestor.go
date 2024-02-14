package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	mt "github.com/txaty/go-merkletree"
)

var CredentialExpiresAt = "ExprTime"
var Values = "Values"

type leaf struct {
	data []byte
}

type merkleTreeLeaf struct {
	VCExpiration time.Time
	Values       []interface{}
}

func (l *leaf) Serialize() ([]byte, error) {
	return l.data, nil
}

// abiEncode mirrors solidity abi encoding
// converts data struct into binary format used by openzep merkle tree library
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
	config    mt.Config
	log       *zerolog.Logger
	abi       *abi.ABI
	producer  sarama.SyncProducer
	schema    [32]byte
	recipient common.Address
	contract  common.Address
	mtxTopic  string
}

// NewAttestor creates a new merkle tree attestation service
// for OpenZep compatilibity:
// OZ sorts all abi encoded leaves before beginning to build tree
// as yet, there is no parameter to make this library sort
// the leaves a priori (sort sibling pairs true refers
// to sorting hashed sib pairs in later steps)
// thus, abi encoding and serialization is performed
// as a first step before we pass any data to the merkle tree func
func NewAttestor(producer sarama.SyncProducer, settings *config.Settings, logger *zerolog.Logger) (*Attestor, error) {
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
		config: mt.Config{
			SortSiblingPairs: true, // parameter for OpenZeppelin compatibility
			HashFunc: func(data []byte) ([]byte, error) {
				return crypto.Keccak256(data), nil
			},
			Mode:               mt.ModeProofGenAndTreeBuild,
			DisableLeafHashing: true, // also required for OpenZeppelin compatilibity
		},
		log:       logger,
		abi:       abi,
		producer:  producer,
		schema:    byteArray,
		contract:  common.HexToAddress(settings.AttestationContract),
		recipient: common.HexToAddress(settings.AttestationRecipient),
		mtxTopic:  settings.MetaTransactionSendTopic,
	}, nil
}

func (a *Attestor) GenerateMerkleTree(data map[string]merkleTreeLeaf, dataTypes []string) (*mt.MerkleTree, error) {
	encodedLeaves := [][]byte{}
	for _, device := range data {
		codeBytes, err := abiEncode(dataTypes, device.Values)
		if err != nil {
			a.log.Err(err).Msg("failed to abi encode attestation data")
		}
		encodedLeaves = append(encodedLeaves, crypto.Keccak256(crypto.Keccak256(codeBytes)))
	}

	// sort abi encoded values for OpenZeppelin compatibility
	sort.Slice(encodedLeaves, func(i, j int) bool {
		return bytes.Compare(encodedLeaves[i], encodedLeaves[j]) < 0
	})

	blocks := []mt.DataBlock{}
	for _, vc := range encodedLeaves {
		blocks = append(blocks, &leaf{data: vc})
	}

	return mt.New(&a.config, blocks)
}

func (a *Attestor) GenerateAttestation(root [32]byte) ([]byte, error) {
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

func (a *Attestor) AttestOnChain(data []byte) (string, error) {
	reqID := ksuid.New().String()
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
			Data: data,
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
	// TODO ae: use reqID when consuming claim to make sure it goes through
	return reqID, err
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
