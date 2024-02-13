package services

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
	config   mt.Config
	log      *zerolog.Logger
	abi      *abi.ABI
	producer sarama.SyncProducer
}

func NewAttestor(producer sarama.SyncProducer, logger *zerolog.Logger) (*Attestor, error) {
	abi, err := contracts.EasMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &Attestor{
		config: mt.Config{
			SortSiblingPairs: true, // parameter for OpenZeppelin compatibility
			HashFunc: func(data []byte) ([]byte, error) {
				return crypto.Keccak256(data), nil
			},
			Mode:               mt.ModeProofGenAndTreeBuild,
			DisableLeafHashing: true, // also required for OpenZeppelin compatilibity
		},
		log:      logger,
		abi:      abi,
		producer: producer,
	}, nil
}

func (a *Attestor) GenerateMerkleTree(data map[string]map[string]interface{}, dataTypes []string) (*mt.MerkleTree, error) {
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

func (a *Attestor) GenerateAttestation(root []byte) ([]byte, error) {
	args := abi.Arguments{}
	// TODO ae: this should be bytes
	abiType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	args = append(args, abi.Argument{Type: abiType})

	enocodedData, err := args.Pack(common.Bytes2Hex(root)) // TODO ae: double check this work, was using dummy data before
	if err != nil {
		return nil, err
	}

	var byteArray [32]byte

	// TODO ae: get this from settings and initialize when we make the attestor
	schemaSlice := common.Hex2Bytes("93dbaf80a47c49a96e9ad6e038a064088d322f9b42d4e4bd8e78efb947b448ad")
	for i := 0; i < len(schemaSlice) && i < len(byteArray); i++ {
		byteArray[i] = schemaSlice[i]
	}

	// TODO ae: who do we want the recipient to be?
	req := contracts.AttestationRequest{
		Schema: byteArray,
		Data: contracts.AttestationRequestData{
			Recipient: common.HexToAddress("Be396b4B84a4139EA5C4dbaDde98AE8364eB21e8"),
			Data:      enocodedData,
			Value:     big.NewInt(0),
		},
	}
	return a.abi.Pack("attest", req)
}

func (a *Attestor) AttestOnChain(data []byte) (string, error) {
	reqID := ksuid.New().String()
	// TODO ae: get topic from settings
	_, _, err := a.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: "topic.transaction.request.send",
			Key:   sarama.StringEncoder(reqID),
			Value: sarama.ByteEncoder(data),
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
