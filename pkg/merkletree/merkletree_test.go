package merkletree

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixtureLeaf struct {
	Account string `json:"account"`
	Amount  string `json:"amount"`
}

type fixtureCase struct {
	Distributor string              `json:"distributor"`
	PoolID      string              `json:"poolId"`
	Week        string              `json:"week"`
	Leaves      []fixtureLeaf       `json:"leaves"`
	Root        string              `json:"root"`
	Proofs      map[string][]string `json:"proofs"`
}

func loadFixtures(t *testing.T) []fixtureCase {
	t.Helper()
	b, err := os.ReadFile("testdata/merkle-fixtures.json")
	require.NoError(t, err)
	var cases []fixtureCase
	require.NoError(t, json.Unmarshal(b, &cases))
	require.NotEmpty(t, cases)
	return cases
}

func buildFixtureTree(t *testing.T, fc fixtureCase) *Tree {
	t.Helper()
	poolID, ok := new(big.Int).SetString(fc.PoolID, 10)
	require.True(t, ok)
	week, ok := new(big.Int).SetString(fc.Week, 10)
	require.True(t, ok)
	leaves := make([]Leaf, len(fc.Leaves))
	for i, fl := range fc.Leaves {
		amount, ok := new(big.Int).SetString(fl.Amount, 10)
		require.True(t, ok, "bad amount %q", fl.Amount)
		leaves[i] = Leaf{Account: common.HexToAddress(fl.Account), Amount: amount}
	}
	tree, err := New(common.HexToAddress(fc.Distributor), poolID, week, leaves)
	require.NoError(t, err)
	return tree
}

func TestGoldenRoots(t *testing.T) {
	for i, fc := range loadFixtures(t) {
		t.Run(fmt.Sprintf("case%d_%dleaves", i, len(fc.Leaves)), func(t *testing.T) {
			tree := buildFixtureTree(t, fc)
			assert.Equal(t, strings.ToLower(fc.Root), tree.Root().Hex())
		})
	}
}

func TestGoldenProofs(t *testing.T) {
	for i, fc := range loadFixtures(t) {
		t.Run(fmt.Sprintf("case%d_%dleaves", i, len(fc.Leaves)), func(t *testing.T) {
			tree := buildFixtureTree(t, fc)

			// Check every account for small cases; sample at least 200 for the
			// 5000-leaf case to keep the test fast.
			stride := 1
			if len(fc.Leaves) > 1000 {
				stride = len(fc.Leaves) / 200
			}

			checked := 0
			for i := 0; i < len(fc.Leaves); i += stride {
				account := fc.Leaves[i].Account
				want, ok := fc.Proofs[account]
				require.True(t, ok, "fixture missing proof for %s", account)

				proof, err := tree.Proof(common.HexToAddress(account))
				require.NoError(t, err)

				got := make([]string, len(proof))
				for j, p := range proof {
					got[j] = p.Hex()
				}
				wantLower := make([]string, len(want))
				for j, p := range want {
					wantLower[j] = strings.ToLower(p)
				}
				require.Equal(t, wantLower, got, "proof mismatch for %s", account)
				checked++
			}
			if len(fc.Leaves) > 1000 {
				require.GreaterOrEqual(t, checked, 200)
			}
		})
	}
}

func TestProofUnknownAccount(t *testing.T) {
	fc := loadFixtures(t)[0]
	tree := buildFixtureTree(t, fc)
	_, err := tree.Proof(common.HexToAddress("0x000000000000000000000000000000000000dEaD"))
	require.Error(t, err)
}

func TestNewValidation(t *testing.T) {
	distributor := common.HexToAddress("0x4a679253410272dd5232B3Ff7cF5dbB88f295319")
	poolID, week := big.NewInt(0), big.NewInt(1)
	account := common.HexToAddress("0x04Dc5Cb9DF85a4F71306919EE6Fc38b12AEC9993")
	other := common.HexToAddress("0xa368E3DcEeA4Fa5f98C203149166d3FeD8a38671")

	t.Run("empty leaves", func(t *testing.T) {
		_, err := New(distributor, poolID, week, nil)
		require.Error(t, err)
	})

	t.Run("duplicate account", func(t *testing.T) {
		_, err := New(distributor, poolID, week, []Leaf{
			{Account: account, Amount: big.NewInt(1)},
			{Account: other, Amount: big.NewInt(2)},
			{Account: account, Amount: big.NewInt(3)},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate")
	})

	t.Run("zero amount", func(t *testing.T) {
		_, err := New(distributor, poolID, week, []Leaf{{Account: account, Amount: big.NewInt(0)}})
		require.Error(t, err)
	})

	t.Run("nil amount", func(t *testing.T) {
		_, err := New(distributor, poolID, week, []Leaf{{Account: account, Amount: nil}})
		require.Error(t, err)
	})

	t.Run("negative amount", func(t *testing.T) {
		_, err := New(distributor, poolID, week, []Leaf{{Account: account, Amount: big.NewInt(-1)}})
		require.Error(t, err)
	})
}

func TestSingleLeafRootEqualsLeafHash(t *testing.T) {
	fc := loadFixtures(t)[0]
	require.Len(t, fc.Leaves, 1)
	tree := buildFixtureTree(t, fc)
	proof, err := tree.Proof(common.HexToAddress(fc.Leaves[0].Account))
	require.NoError(t, err)
	require.Empty(t, proof)
	require.Equal(t, strings.ToLower(fc.Root), tree.Root().Hex())
}

func TestMarshalRoundtripAndVerifyRoot(t *testing.T) {
	for _, fc := range loadFixtures(t) {
		if len(fc.Leaves) > 200 {
			continue
		}
		t.Run(fmt.Sprintf("%d_leaves", len(fc.Leaves)), func(t *testing.T) {
			tree := buildFixtureTree(t, fc)
			b, err := json.Marshal(tree)
			require.NoError(t, err)

			tf, err := UnmarshalTreeFile(b)
			require.NoError(t, err)
			require.Equal(t, "dimo-merkle-v1", tf.Format)
			require.Equal(t, tree.Root(), tf.Root)
			require.Len(t, tf.Leaves, len(fc.Leaves))
			require.NoError(t, tf.VerifyRoot())

			// Tamper with one amount: VerifyRoot must fail.
			tf.Leaves[0].Amount = new(big.Int).Add(tf.Leaves[0].Amount, big.NewInt(1))
			require.Error(t, tf.VerifyRoot())
		})
	}
}

func TestMarshalJSONAmountsAreDecimalStrings(t *testing.T) {
	distributor := common.HexToAddress("0x4a679253410272dd5232B3Ff7cF5dbB88f295319")
	// Larger than 2^53; would lose precision as a JSON number.
	amount, ok := new(big.Int).SetString("14250000000000000001", 10)
	require.True(t, ok)
	account := common.HexToAddress("0x04Dc5Cb9DF85a4F71306919EE6Fc38b12AEC9993")

	tree, err := New(distributor, big.NewInt(0), big.NewInt(226), []Leaf{{Account: account, Amount: amount}})
	require.NoError(t, err)

	b, err := json.Marshal(tree)
	require.NoError(t, err)
	require.Contains(t, string(b), `"amount":"14250000000000000001"`)
	require.Contains(t, string(b), `"poolId":"0"`)
	require.Contains(t, string(b), `"week":"226"`)
	require.Contains(t, string(b), `"format":"dimo-merkle-v1"`)
	// Account must be EIP-55 checksummed.
	require.Contains(t, string(b), `"account":"0x04Dc5Cb9DF85a4F71306919EE6Fc38b12AEC9993"`)

	tf, err := UnmarshalTreeFile(b)
	require.NoError(t, err)
	require.Zero(t, tf.Leaves[0].Amount.Cmp(amount), "amount must roundtrip without precision loss")
}

func TestLeavesCanonicalOrder(t *testing.T) {
	fc := loadFixtures(t)[2] // 3-leaf case
	tree := buildFixtureTree(t, fc)
	leaves := tree.Leaves()
	require.Len(t, leaves, len(fc.Leaves))
	// All input accounts present.
	seen := make(map[common.Address]bool, len(leaves))
	for _, l := range leaves {
		seen[l.Account] = true
	}
	for _, fl := range fc.Leaves {
		require.True(t, seen[common.HexToAddress(fl.Account)])
	}
	// Calling Leaves twice returns the same order.
	require.Equal(t, leaves, tree.Leaves())
}

func TestUnmarshalTreeFileErrors(t *testing.T) {
	t.Run("bad json", func(t *testing.T) {
		_, err := UnmarshalTreeFile([]byte(`{`))
		require.Error(t, err)
	})
	t.Run("wrong format", func(t *testing.T) {
		_, err := UnmarshalTreeFile([]byte(`{"format":"standard-v1"}`))
		require.Error(t, err)
	})
	t.Run("negative amount", func(t *testing.T) {
		_, err := UnmarshalTreeFile([]byte(`{
			"format": "dimo-merkle-v1",
			"distributor": "0x4a679253410272dd5232B3Ff7cF5dbB88f295319",
			"poolId": "0",
			"week": "1",
			"root": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"leaves": [
				{"account": "0x04Dc5Cb9DF85a4F71306919EE6Fc38b12AEC9993", "amount": "-5", "proof": []}
			]
		}`))
		require.Error(t, err)
		require.Contains(t, err.Error(), "amount in leaf 0 must be positive")
	})
}
