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

// GetUserRewards godoc
// @Description  A summary of the user's rewards.
// @Success      200 {object} controllers.UserResponse
// @Security     BearerAuth
// @Router       /rewards [get]
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
	// Points is the user's total number of points, across all devices and issuance weeks.
	Points int `json:"points" example:"5000"`
	// Devices is a list of the user's devices, together with some information about their
	// connectivity.
	Devices []*UserResponseDevice `json:"devices"`
	// ThisWeek describes the current issuance week.
	ThisWeek UserResponseThisWeek `json:"thisWeek"`
}

type UserResponseDevice struct {
	// ID is the user device ID used across all services.
	ID string `json:"id" example:"27cv7gVTh9h4RJuTsmJHpBcr4I9"`
	// Points is the total number of points that the device has earned across all weeks.
	Points int `json:"points" example:"5000"`
	// ConnectedThisWeek is true if we've seen activity from the device during the current issuance
	// week.
	ConnectedThisWeek bool `json:"connectedThisWeek" example:"true"`
	// ConnectionStreak is what we consider the streak of the device to be. This may not literally
	// be the number of consecutive connected weeks, because the user may disconnect for a week
	// without penalty, or have the connection streak reduced after three weeks of inactivity.
	ConnectionStreak int `json:"connectionStreak" example:"4"`
	// DisconnectionStreak is the number of consecutive issuance weeks that the device has been
	// disconnected. This number resets to 0 as soon as a device earns rewards for a certain week.
	DisconnectionStreak int `json:"disconnectionStreak,omitempty" example:"0"`
	// Level is the level 1-4 of the device. This is fully determined by ConnectionStreak.
	Level int `json:"level" example:"2"`
}

type UserResponseThisWeek struct {
	// Start is the timestamp of the start of the issuance week.
	Start time.Time `json:"start" example:"2022-04-11T05:00:00Z"`
	// End is the timestamp of the start of the next issuance week.
	End time.Time `json:"end" example:"2022-04-18T05:00:00Z"`
}
