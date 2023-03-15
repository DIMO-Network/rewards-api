package controllers

import (
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

// GetBalanceHistory godoc
// @Description  A summary of the user's DIMO balance across all chains.
// @Success      200 {object} controllers.BalanceHistory
// @Security     BearerAuth
// @Router       /user/history/balance [get]
func (r *RewardsController) GetBalanceHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	user, err := r.UsersClient.GetUser(c.Context(), &pb_users.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user data.")
		return opaqueInternalError
	}

	balanceHistory := BalanceHistory{
		BalanceHistory: []Balance{},
	}

	if user.EthereumAddress == nil {
		return c.JSON(balanceHistory)
	}

	addr := common.HexToAddress(*user.EthereumAddress)

	// Terrible no good.
	tfs, err := models.TokenTransfers(
		qm.Where(models.TokenTransferTableColumns.AddressFrom+" != "+models.TokenTransferTableColumns.AddressTo),
		models.TokenTransferWhere.Amount.NEQ(types.NewDecimal(decimal.New(0, 0))),
		qm.Expr(
			models.TokenTransferWhere.AddressTo.EQ(addr.Bytes()),
			qm.Or2(models.TokenTransferWhere.AddressFrom.EQ(addr.Bytes())),
		),
		qm.OrderBy(models.TokenTransferTableColumns.BlockTimestamp+" ASC"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		return err
	}

	runningBalance := big.NewInt(0)

	for _, tf := range tfs {
		delta := tf.Amount.Int(tf.Amount.Int(nil))
		if common.BytesToAddress(tf.AddressFrom) == addr {
			delta = new(big.Int).Mul(delta, big.NewInt(-1))
		}

		runningBalance = new(big.Int).Add(runningBalance, delta)

		if len(balanceHistory.BalanceHistory) == 0 || balanceHistory.BalanceHistory[len(balanceHistory.BalanceHistory)-1].Time != tf.BlockTimestamp {
			balanceHistory.BalanceHistory = append(balanceHistory.BalanceHistory, Balance{Time: tf.BlockTimestamp, Balance: runningBalance})
		} else {
			balanceHistory.BalanceHistory[len(balanceHistory.BalanceHistory)-1].Balance = new(big.Int).Add(balanceHistory.BalanceHistory[len(balanceHistory.BalanceHistory)-1].Balance, delta)
		}
	}

	return c.JSON(balanceHistory)
}

type BalanceHistory struct {
	BalanceHistory []Balance `json:"balanceHistory"`
}

type Balance struct {
	// Time is the block timestamp of this balance update.
	Time time.Time `json:"time" swaggertype:"string" example:"2023-03-06T09:11:00Z"`
	// Balance is the total amount of $DIMO held at this time, across all chains.
	Balance *big.Int `json:"balance" swaggertype:"number" example:"237277217092548851191"`
}
