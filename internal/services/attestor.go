package services

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
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
	config mt.Config
	log    *zerolog.Logger
}

func NewAttestor(logger *zerolog.Logger) *Attestor {
	return &Attestor{
		config: mt.Config{
			SortSiblingPairs: true, // parameter for OpenZeppelin compatibility
			HashFunc: func(data []byte) ([]byte, error) {
				return crypto.Keccak256(data), nil
			},
			Mode:               mt.ModeProofGenAndTreeBuild,
			DisableLeafHashing: true, // also required for OpenZeppelin compatilibity
		},
		log: logger,
	}
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

func (a *Attestor) VerifyAttestation(root, node []byte, proof [][]byte) (bool, error) {
	dataBlock := &leaf{data: node}
	p := mt.Proof{
		Siblings: proof,
	}
	return mt.Verify(dataBlock, &p, root, &a.config)

}
