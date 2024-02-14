package services

import (
	"bytes"
	"fmt"
	"sort"
	"time"

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
	config mt.Config
	log    *zerolog.Logger
}

// NewAttestor creates a new merkle tree attestation service
// for OpenZep compatilibity:
// OZ sorts all abi encoded leaves before beginning to build tree
// as yet, there is no parameter to make this library sort
// the leaves a priori (sort sibling pairs true refers
// to sorting hashed sib pairs in later steps)
// thus, abi encoding and serialization is performed
// as a first step before we pass any data to the merkle tree func
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

func (a *Attestor) VerifyAttestation(root, node []byte, proof [][]byte) (bool, error) {
	dataBlock := &leaf{data: node}
	p := mt.Proof{
		Siblings: proof,
	}
	return mt.Verify(dataBlock, &p, root, &a.config)

}
