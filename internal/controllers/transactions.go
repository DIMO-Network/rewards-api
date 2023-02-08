package controllers

import (
	"errors"
	"math/big"

	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetTransactionHistory godoc
// @Description  A summary of the user's DIMO transaction history, all time.
// @Success      200 {object} controllers.HistoryResponse
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

	incoming, err := models.TokenTransfers(
		models.TokenTransferWhere.UserAddressTo.EQ([]byte(*user.EthereumAddress)),
		qm.OrderBy(models.TokenTransferColumns.CreatedAt+" asc"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving incoming transactions.")
		return opaqueInternalError
	}

	outgoing, err := models.TokenTransfers(
		models.TokenTransferWhere.UserAddressTo.EQ([]byte(*user.EthereumAddress)),
		qm.OrderBy(models.TokenTransferColumns.CreatedAt+" asc"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving outgoing transactions.")
		return opaqueInternalError
	}

	var txHistory TransactionHistory
	if len(incoming) == 0 && len(outgoing) == 0 {
		return c.JSON(txHistory)
	}

	txHistory.IncomingTransaction = make([]IncomingTransactionResponse, len(incoming))
	for n, tx := range incoming {
		amnt, ok := tx.Amount.Int64()
		if !ok {
			logger.Err(errors.New("unable to read tx amount")).Msg("error reading tx amount")
		}

		txHistory.IncomingTransaction[n].Amount = big.NewInt(amnt)
		txHistory.IncomingTransaction[n].FromAddress = common.BytesToAddress(tx.UserAddressFrom)
		txHistory.IncomingTransaction[n].ToAddress = common.BytesToAddress(tx.UserAddressTo)
	}

	txHistory.OutgoingTransactions = make([]OutgoingTransactionResponse, len(outgoing))
	for n, tx := range outgoing {
		amnt, ok := tx.Amount.Int64()
		if !ok {
			logger.Err(errors.New("unable to read tx amount")).Msg("error reading tx amount")
		}
		txHistory.OutgoingTransactions[n].Amount = big.NewInt(amnt)
		txHistory.OutgoingTransactions[n].FromAddress = common.BytesToAddress(tx.UserAddressFrom)
		txHistory.OutgoingTransactions[n].ToAddress = common.BytesToAddress(tx.UserAddressTo)
	}

	return c.JSON(txHistory)
}

type TransactionHistory struct {
	OutgoingTransactions []OutgoingTransactionResponse `json:"outgoingTransactions"`
	IncomingTransaction  []IncomingTransactionResponse `json:"incomingTransactions"`
}

type OutgoingTransactionResponse struct {
	Time        string         `json:"time"`
	ToAddress   common.Address `json:"to"`
	FromAddress common.Address `json:"from"`
	Amount      *big.Int       `json:"amount"`
}

type IncomingTransactionResponse struct {
	Time        string         `json:"time"`
	ToAddress   common.Address `json:"to"`
	FromAddress common.Address `json:"from"`
	Amount      *big.Int       `json:"amount"`
}
