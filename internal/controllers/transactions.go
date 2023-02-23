package controllers

import (
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetTransactionHistory godoc
// @Description  A summary of the user's DIMO transaction history, all time.
// @Success      200 {object} controllers.TransactionHistory
// @Security     BearerAuth
// @Router       /user/history/transactions [get]
func (r *RewardsController) GetTransactionHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	user, err := r.UsersClient.GetUser(c.Context(), &pb_users.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user data.")
		return opaqueInternalError
	}

	txHistory := TransactionHistory{
		Transactions: []APITransaction{},
	}

	if user.EthereumAddress == nil {
		return c.JSON(txHistory)
	}

	addr := common.HexToAddress(*user.EthereumAddress)

	txes, err := models.TokenTransfers(
		models.TokenTransferWhere.AddressTo.EQ(addr.Bytes()),
		qm.Or2(models.TokenTransferWhere.AddressFrom.EQ(addr.Bytes())),
		qm.OrderBy(models.TokenTransferColumns.BlockTimestamp+" DESC"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving incoming transactions.")
		return opaqueInternalError
	}

	for _, tx := range txes {
		apiTx := APITransaction{
			ChainID: tx.ChainID,
			Time:    tx.BlockTimestamp,
			From:    common.BytesToAddress(tx.AddressFrom),
			To:      common.BytesToAddress(tx.AddressTo),
			Value:   tx.Amount.Int(nil),
		}
		txHistory.Transactions = append(txHistory.Transactions, apiTx)
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
}
