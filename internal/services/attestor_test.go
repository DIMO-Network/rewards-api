package services

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMerkleTreeGeneration(t *testing.T) {
	assert := assert.New(t)
	logger := zerolog.Nop()
	witness := NewAttestor(&logger)

	attestationData := map[string]merkleTreeLeaf{}
	for i := 1; i < 5; i++ {
		dummyVin := fmt.Sprintf("vin%d", i)
		currentVCExpiration := time.Now()
		potentialEarner, ok := attestationData[fmt.Sprintf("vin%d", i)]
		if ok {
			if potentialEarner.VCExpiration.Before(currentVCExpiration) {
				attestationData[dummyVin] = merkleTreeLeaf{
					VCExpiration: currentVCExpiration,
					Values:       []interface{}{uint64(i), "deviceDefinitionId", "latestVinCredentialId", "signature"},
				}
			}
		} else {
			attestationData[dummyVin] = merkleTreeLeaf{
				VCExpiration: currentVCExpiration,
				Values:       []interface{}{uint64(i), "deviceDefinitionId", "latestVinCredentialId", "signature"},
			}
		}
	}

	sampleTypes := []string{"uint64", "string", "string", "string"}
	tree, err := witness.GenerateMerkleTree(attestationData, sampleTypes)
	assert.NoError(err)

	encodedLeaves := [][]byte{}
	for _, device := range attestationData {
		codeBytes, err := abiEncode(sampleTypes, device.Values)
		assert.NoError(err)
		encodedLeaves = append(encodedLeaves, crypto.Keccak256(crypto.Keccak256(codeBytes)))
	}

	// sort values for OpenZeppelin compatibility
	sort.Slice(encodedLeaves, func(i, j int) bool {
		return bytes.Compare(encodedLeaves[i], encodedLeaves[j]) < 0
	})

	assert.Equal(len(tree.Proofs), len(encodedLeaves))
	for i, proof := range tree.Proofs {
		valid, err := witness.VerifyAttestation(tree.Root, encodedLeaves[i], proof.Siblings)
		assert.NoError(err)
		assert.True(valid)
	}
}
