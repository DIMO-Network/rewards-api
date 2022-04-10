package controllers

import (
	"time"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/devices"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type RewardsController struct {
	DB            func() *database.DBReaderWriter
	Logger        *zerolog.Logger
	DataClient    services.DeviceDataClient
	IntegClient   pb.IntegrationServiceClient
	DevicesClient pb.UserDeviceServiceClient
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

var opaqueInternalError = fiber.NewError(fiber.StatusInternalServerError, "Internal error.")

func (r *RewardsController) GetUserRewards(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	now := time.Now()
	weekNum := services.GetWeekNum(now)
	weekStart := services.NumToWeekStart(weekNum)

	devices, err := r.DevicesClient.ListUserDevicesForUser(c.Context(), &pb.ListUserDevicesForUserRequest{
		UserId: userID,
	})
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user's devices.")
		return opaqueInternalError
	}

	outLi := make([]*UserResponseDevice, len(devices.UserDevices))
	userPts := 0

	for i, device := range devices.UserDevices {
		dlog := logger.With().Str("userDeviceId", device.Id).Logger()
		lastActive, seen, err := r.DataClient.GetLastActivity(device.Id)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve last activity.")
			return opaqueInternalError
		}
		var activeThisWeek = false
		if seen && !lastActive.Before(weekStart) {
			activeThisWeek = true
		}
		rewards, err := models.Rewards(
			models.RewardWhere.UserDeviceID.EQ(device.Id),
			models.RewardWhere.UserID.EQ(userID),
			qm.OrderBy(models.RewardColumns.IssuanceWeekID+" desc"),
		).All(c.Context(), r.DB().Reader)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve previously earned rewards.")
			return opaqueInternalError
		}

		pts := 0
		for _, r := range rewards {
			pts += r.StreakPoints + r.IntegrationPoints
		}

		userPts += pts

		lvl := 1
		connectionStreak := 0
		disconnectionStreak := 0
		if len(rewards) > 0 {
			lvl = services.GetLevel(rewards[0].ConnectionStreak).Level
			connectionStreak = rewards[0].ConnectionStreak
			disconnectionStreak = rewards[0].DisconnectionStreak
		}

		outLi[i] = &UserResponseDevice{
			ID:                  device.Id,
			Points:              pts,
			ConnectedThisWeek:   activeThisWeek,
			ConnectionStreak:    connectionStreak,
			DisconnectionStreak: disconnectionStreak,
			Level:               lvl,
		}
	}

	return c.JSON(UserResponse{
		Points: userPts,
		ThisWeek: UserResponseThisWeek{
			Start: weekStart,
			End:   services.NumToWeekEnd(weekNum),
		},
		Devices: outLi,
	})
}

type UserResponse struct {
	Points   int                   `json:"points"`
	Devices  []*UserResponseDevice `json:"devices"`
	ThisWeek UserResponseThisWeek  `json:"thisWeek"`
}

type UserResponseDevice struct {
	ID                  string `json:"id"`
	Points              int    `json:"points"`
	ConnectedThisWeek   bool   `json:"connectedThisWeek"`
	ConnectionStreak    int    `json:"connectionStreak"`
	DisconnectionStreak int    `json:"disconnectionStreak,omitempty"`
	Level               int    `json:"level"`
}

type UserResponseThisWeek struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
