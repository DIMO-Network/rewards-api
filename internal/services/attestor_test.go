package services

import (
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

	attestationData := make(map[string]map[string]interface{}, 0)
	for i := 1; i < 5; i++ {
		dummyVin := fmt.Sprintf("vin%d", i)
		currentVCExpiration := time.Now()
		earningVehicle, ok := attestationData[fmt.Sprintf("vin%d", i)]
		if ok {
			credExp, ok := earningVehicle[CredentialExpiresAt]
			if ok {
				if credExp.(time.Time).Before(currentVCExpiration) {
					attestationData[dummyVin] = map[string]interface{}{
						CredentialExpiresAt: currentVCExpiration,
						Values:              []interface{}{uint64(i), "deviceDefinitionId", "latestVinCredentialId", "signature"},
					}
				}
			}
		}
		attestationData[dummyVin] = map[string]interface{}{
			CredentialExpiresAt: currentVCExpiration,
			Values:              []interface{}{uint64(i), "deviceDefinitionId", "latestVinCredentialId", "signature"},
		}
	}

	sampleTypes := []string{"uint64", "string", "string", "string"}
	tree, err := witness.GenerateMerkleTree(attestationData, sampleTypes)
	assert.NoError(err)

	encodedLeaves := [][]byte{}
	for _, device := range attestationData {
		codeBytes, err := abiEncode(sampleTypes, device[Values].([]interface{}))
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
