package controllers

import (
	"math/big"
	"strconv"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
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
	maybeAddr, err := r.getCallerEthAddress(c)
	if err != nil {
		return err
	}

	balanceHistory := BalanceHistory{
		BalanceHistory: []Balance{},
	}

	if maybeAddr == nil {
		return c.JSON(balanceHistory)
	}

	addr := *maybeAddr

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
		value := tf.Amount.Int(nil)
		if common.BytesToAddress(tf.AddressFrom) == addr {
			runningBalance = new(big.Int).Sub(runningBalance, value)
		} else {
			runningBalance = new(big.Int).Add(runningBalance, value)
		}

		if l := len(balanceHistory.BalanceHistory); l == 0 || balanceHistory.BalanceHistory[l-1].Time != tf.BlockTimestamp {
			balanceHistory.BalanceHistory = append(balanceHistory.BalanceHistory, Balance{Time: tf.BlockTimestamp, Balance: runningBalance})
		} else {
			balanceHistory.BalanceHistory[l-1].Balance = runningBalance
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

// GetPotentialTokens godoc
// @Description Calculate potential DIMO token earnings for a given week and points
// @Param        date   query    string  true  "Week date (YYYY-MM-DD)"
// @Param        points query    int     true  "Number of points"
// @Success      200    {object} PotentialTokensResponse
// @Router       /rewards/potential [get]
func (r *RewardsController) GetPotentialTokens(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	date := time.Now().Add(-7 * 24 * time.Hour)
	if dateStr != "" {
		var err error
		date, err = time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		}
	}

	points, err := strconv.Atoi(c.Query("points"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid points value")
	}

	dbStorage := storage.DBStorage{DBS: r.DB, Logger: r.Logger}
	potentialTokens, err := dbStorage.CalculateTokensForPoints(c.Context(), points, date)
	if err != nil {
	}

	return c.JSON(PotentialTokensResponse{
		Date:        date,
		Points:      points,
		TokenAmount: potentialTokens,
	})
}

type PotentialTokensResponse struct {
	Date        time.Time    `json:"date"`
	Points      int          `json:"points"`
	TokenAmount *decimal.Big `json:"tokenAmount"`
}
