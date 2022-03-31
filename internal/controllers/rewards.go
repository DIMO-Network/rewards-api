package controllers

import (
	"time"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
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

func (r *RewardsController) AdminGetRewards(c *fiber.Ctx) error {
	userID := c.Params("userID")
	rewards, err := models.Rewards(
		models.RewardWhere.UserID.EQ(userID),
		qm.OrderBy(
			models.RewardColumns.IssuanceWeekID+" desc, "+models.RewardColumns.UserDeviceID+" asc",
		),
	).All(c.Context(), r.DB().Reader)
	if err != nil {
		r.Logger.Err(err).Msg("")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}

	weeks := []*IssuanceWeekResponse{}
	var week *IssuanceWeekResponse

	points := 0

	issuanceWeek := -1
	for _, reward := range rewards {
		if reward.IssuanceWeekID != issuanceWeek {
			week = &IssuanceWeekResponse{
				Start: services.NumToWeekStart(reward.IssuanceWeekID),
				End:   services.NumToWeekEnd(reward.IssuanceWeekID),
			}
			weeks = append(weeks, week)
			issuanceWeek = reward.IssuanceWeekID
		}
		week.Devices = append(week.Devices, &IssuanceWeekDeviceResponse{
			DeviceID:                  reward.UserDeviceID,
			EffectiveConnectionStreak: reward.EffectiveConnectionStreak,
			DisconnectionStreak:       reward.DisconnectionStreak,
			StreakPoints:              reward.StreakPoints,
			IntegrationIDs:            reward.IntegrationIds,
			IntegrationPoints:         reward.IntegrationPoints,
			Points:                    reward.StreakPoints + reward.IntegrationPoints,
		})
		week.Points += reward.StreakPoints + reward.IntegrationPoints
		points += reward.StreakPoints + reward.IntegrationPoints
	}

	resp := UserResponse{
		UserID:        userID,
		IssuanceWeeks: weeks,
		Points:        points,
	}

	return c.JSON(resp)
}

type UserResponse struct {
	UserID        string                  `json:"userId"`
	Points        int                     `json:"points"`
	IssuanceWeeks []*IssuanceWeekResponse `json:"issuanceWeeks"`
}

type IssuanceWeekResponse struct {
	Start   time.Time                     `json:"start"`
	End     time.Time                     `json:"end"`
	Points  int                           `json:"points"`
	Devices []*IssuanceWeekDeviceResponse `json:"devices"`
}

type IssuanceWeekDeviceResponse struct {
	DeviceID                  string   `json:"deviceId"`
	Points                    int      `json:"points"`
	EffectiveConnectionStreak int      `json:"effectiveConnectionStreak"`
	DisconnectionStreak       int      `json:"disconnectionStreak"`
	StreakPoints              int      `json:"streakPoints"`
	IntegrationIDs            []string `json:"integrationIds"`
	IntegrationPoints         int      `json:"integrationPoints"`
}
