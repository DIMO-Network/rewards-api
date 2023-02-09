package controllers

import (
	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
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
		models.TokenTransferWhere.AddressTo.EQ([]byte(*user.EthereumAddress)),
		qm.OrderBy(models.TokenTransferColumns.BlockTimestamp+" asc"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving incoming transactions.")
		return opaqueInternalError
	}

	outgoing, err := models.TokenTransfers(
		models.TokenTransferWhere.AddressTo.EQ([]byte(*user.EthereumAddress)),
		qm.OrderBy(models.TokenTransferColumns.BlockTimestamp+" asc"),
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
		txHistory.IncomingTransaction[n].Amount = tx.Amount
		txHistory.IncomingTransaction[n].FromAddress = common.BytesToAddress(tx.AddressFrom)
		txHistory.IncomingTransaction[n].ToAddress = common.BytesToAddress(tx.AddressTo)
		txHistory.IncomingTransaction[n].Time = tx.BlockTimestamp.String()
	}

	txHistory.OutgoingTransactions = make([]OutgoingTransactionResponse, len(outgoing))
	for n, tx := range outgoing {
		txHistory.OutgoingTransactions[n].Amount = tx.Amount
		txHistory.OutgoingTransactions[n].FromAddress = common.BytesToAddress(tx.AddressFrom)
		txHistory.OutgoingTransactions[n].ToAddress = common.BytesToAddress(tx.AddressTo)
		txHistory.IncomingTransaction[n].Time = tx.BlockTimestamp.String()
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
	Amount      types.Decimal  `json:"amount"`
}

type IncomingTransactionResponse struct {
	Time        string         `json:"time"`
	ToAddress   common.Address `json:"to"`
	FromAddress common.Address `json:"from"`
	Amount      types.Decimal  `json:"amount"`
}
