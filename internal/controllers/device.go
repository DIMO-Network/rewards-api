package controllers

import (
	"math/big"
	"strconv"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/ericlagergren/decimal"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type DeviceController struct {
	DB     db.Store
	Logger *zerolog.Logger
}

func (r *DeviceController) GetDevice(c *fiber.Ctx) error {
	ts := c.Params("tokenID")

	ti, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldn't parse token id.")
	}

	rs, err := models.Rewards(
		models.RewardWhere.AftermarketTokenID.EQ(types.NewNullDecimal(decimal.New(ti, 0))),
		qm.OrderBy(models.RewardColumns.IssuanceWeekID+" DESC"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		return err
	}

	out := make([]DeviceWeek, len(rs))

	for i, r := range rs {
		out[i] = DeviceWeek{
			Start:  services.NumToWeekStart(r.IssuanceWeekID),
			End:    services.NumToWeekEnd(r.IssuanceWeekID),
			Tokens: r.Tokens.Int(nil),
		}
	}

	return c.JSON(DeviceSummary{Weeks: out})
}

type DeviceWeek struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Tokens *big.Int  `json:"tokens"`
}

type DeviceSummary struct {
	Weeks []DeviceWeek `json:"weeks"`
}
