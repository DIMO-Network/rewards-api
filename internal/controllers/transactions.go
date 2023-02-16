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
			ChainID:     tx.ChainID,
			Time:        tx.BlockTimestamp,
			FromAddress: common.BytesToAddress(tx.AddressFrom),
			ToAddress:   common.BytesToAddress(tx.AddressTo),
			Value:       tx.Amount.Int(nil),
		}
		txHistory.Transactions = append(txHistory.Transactions, apiTx)
	}

	return c.JSON(txHistory)
}

type TransactionHistory struct {
	Transactions []APITransaction `json:"transactions"`
}

type APITransaction struct {
	ChainID     int64          `json:"chainId"`
	Time        time.Time      `json:"time"`
	ToAddress   common.Address `json:"to"`
	FromAddress common.Address `json:"from"`
	Value       *big.Int       `json:"value"`
}
