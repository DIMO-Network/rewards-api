package controllers

import (
	"math/big"
	"strconv"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
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
	addr, err := GetTokenEthAddr(c)
	if err != nil {
		return err
	}

	balanceHistory := BalanceHistory{
		BalanceHistory: []Balance{},
	}

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

// GetHistoricalConversion godoc
// @Description Calculate DIMO token earned fo a given week and popints
// @Param        points query    int     true  "Number of points"
// @Param        time   query    string  false  "Time in the week to calculate potential tokens earned based on the provided points (defaults to last week) (format RFC-3339 e.x. 2024-12-23T12:41:42Z)"
// @Success      200    {object} HistoricalConversionResponse
// @Router       /rewards/convert [get]
func (r *RewardsController) GetHistoricalConversion(c *fiber.Ctx) error {
	dateStr := c.Query("time")
	var weekID int
	if dateStr != "" {
		weekTime, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid time format use RFC-3339 e.g. (2024-12-23T12:41:42Z)")
		}
		weekID = date.GetWeekNum(weekTime)
	} else {
		weekID = date.GetWeekNum(time.Now()) - 1
	}

	points, err := strconv.Atoi(c.Query("points"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "couldn't parse points value as a number")
	}
	if points < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "points must be positive")
	}

	finishedWeekID, potentialTokens, err := storage.CalculateTokensForPoints(c.Context(), r.DB, points, weekID)
	if err != nil {
		return err
	}

	return c.JSON(HistoricalConversionResponse{
		Points:      points,
		Tokens:      potentialTokens,
		StartOfWeek: date.NumToWeekStart(finishedWeekID),
	})
}

type HistoricalConversionResponse struct {
	// Points is the number of points used to calculate the potential tokens.
	Points int `json:"points"`
	// Tokens is the number of tokens ($DIMO/eth not wei) that would be earned for the given number of points.
	Tokens *decimal.Big `json:"tokens"`
	// StartOfWeek is the start of the week for the conversion.
	StartOfWeek time.Time `json:"startOfWeek"`
}
