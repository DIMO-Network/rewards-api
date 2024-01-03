package utils

import (
	"math/big"

	"github.com/volatiletech/sqlboiler/v4/types"
)

func NullDecimalToIntDefaultZero(num types.NullDecimal) *big.Int {
	if num.IsZero() {
		return big.NewInt(0)
	}

	return num.Int(nil)
}
