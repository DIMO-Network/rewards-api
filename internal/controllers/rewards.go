package controllers

import (
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type RewardsController struct {
	DB     func() *database.DBReaderWriter
	Logger *zerolog.Logger
}

type RewardsResponse struct {
	UserID      string `json:"userId"`
	TotalPoints int    `json:"totalPoints"`
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}

func (r *RewardsController) GetRewards(c *fiber.Ctx) error {
	userID := getUserID(c)
	rewards, err := models.Rewards(
		models.RewardWhere.UserID.EQ(userID),
		qm.OrderBy(
			models.RewardColumns.IssuanceWeekID+" desc",
			models.RewardColumns.UserDeviceID+" asc",
		),
	).All(c.Context(), r.DB().Reader)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}
	return c.JSON(rewards)
}
