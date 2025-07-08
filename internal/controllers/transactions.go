package controllers

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

// GetTransactionHistory godoc
// @Description  A summary of the user's DIMO transaction history, all time.
// @Success      200 {object} controllers.TransactionHistory
// @Security     BearerAuth
// @Param        type query string false "A label for a transaction type." Enums(Baseline, Referrals, Marketplace, Other)
// @Router       /user/history/transactions [get]
// @Deprecated
func (r *RewardsController) GetTransactionHistory(c *fiber.Ctx) error {
	txHistory := TransactionHistory{
		Transactions: []APITransaction{},
	}

	return c.JSON(txHistory)
}

type TransactionHistory struct {
	Transactions []APITransaction `json:"transactions"`
}

type APITransaction struct {
	// ChainID is the chain id of the chain on which the transaction took place. Important
	// values are 137 for Polygon, 1 for Ethereum.
	ChainID int64 `json:"chainId" example:"137"`
	// Time is the timestamp of the block in which the transaction took place, in RFC-3999 format.
	Time time.Time `json:"time" example:"2023-01-22T09:00:12Z"`
	// From is the address of the source of the value, in 0x-prefixed hex.
	From common.Address `json:"from" example:"0xf316832fbfe49f90df09eee019c2ece87fad3fac" swaggertype:"string"`
	// To is the address of the recipient of the value, in 0x-prefixed hex.
	To common.Address `json:"to" example:"0xc66d80f5063677425270013136ef9fa2bf1f9f1a" swaggertype:"string"`
	// Value is the amount of token being transferred. Divide by 10^18 to get what people
	// normally consider $DIMO.
	Value *big.Int `json:"value" example:"10000000000000000" swaggertype:"number"`
	// Type is a transaction type.
	Type *string `json:"type,omitempty" enums:"Baseline,Referrals,Marketplace"`
	// Description is a short elaboration of the Type or a generic, e.g., "Incoming transfer" message.
	Description string `json:"description"`
}
