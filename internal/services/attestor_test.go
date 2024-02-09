package services

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMerkleTreeGeneration(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	logger := zerolog.Nop()
	witness := NewAttestor(&logger)

	sampleValues := [][]interface{}{}
	for i := 1; i < 10; i++ {
		usr := []interface{}{}
		usr = append(usr, uint64(i))
		usr = append(usr, "test")
		usr = append(usr, "test")
		usr = append(usr, "test")

		sampleValues = append(sampleValues, usr)
	}

	sampleTypes := []string{"uint64", "string", "string", "string"}
	fmt.Println(sampleValues)
	tree, err := witness.GenerateMerkleTree(ctx, sampleValues, sampleTypes)
	assert.NoError(err)

	encodedLeaves := [][]byte{}
	for _, d := range sampleValues {
		codeBytes, err := abiEncode(sampleTypes, d)
		assert.NoError(err)
		encodedLeaves = append(encodedLeaves, crypto.Keccak256(crypto.Keccak256(codeBytes)))
	}

	// sort values for OpenZeppelin compatibility
	sort.Slice(encodedLeaves, func(i, j int) bool {
		return string(encodedLeaves[i]) < string(encodedLeaves[j])
	})

	assert.Equal(len(tree.Proofs), len(encodedLeaves))
	for i, proof := range tree.Proofs {
		valid, err := witness.VerifyAttestation(tree.Root, encodedLeaves[i], proof.Siblings)
		assert.NoError(err)
		assert.True(valid)
	}
}
