package controllers

import (
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/devices"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RewardsController struct {
	DB            func() *database.DBReaderWriter
	Logger        *zerolog.Logger
	DataClient    services.DeviceDataClient
	IntegClient   pb.IntegrationServiceClient
	DevicesClient pb.UserDeviceServiceClient
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
// @Router       /user [get]
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

	intDescs, err := r.IntegClient.ListIntegrations(c.Context(), &emptypb.Empty{})
	if err != nil {
		return opaqueInternalError
	}

	intMap := make(map[string]string)
	for _, intDesc := range intDescs.Integrations {
		intMap[intDesc.Vendor] = intDesc.Id
	}

	outLi := make([]*UserResponseDevice, len(devices.UserDevices))
	userPts := 0

	for i, device := range devices.UserDevices {
		dlog := logger.With().Str("userDeviceId", device.Id).Logger()
		var maybeLastActive *time.Time
		lastActive, seen, err := r.DataClient.GetLastActivity(device.Id)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve last activity.")
			return opaqueInternalError
		}

		outInts := make([]UserResponseIntegration, 0)

		var activeThisWeek = false
		if seen {
			maybeLastActive = &lastActive
			if !lastActive.Before(weekStart) {
				activeThisWeek = true

				ints, err := r.DataClient.GetIntegrations(device.Id, weekStart, now)
				if err != nil {
					return opaqueInternalError
				}

				if services.ContainsString(ints, intMap["AutoPi"]) {
					outInts = append(outInts, UserResponseIntegration{
						ID:     intMap["AutoPi"],
						Vendor: "AutoPi",
						Points: 6000,
					})
					if services.ContainsString(ints, intMap["SmartCar"]) {
						outInts = append(outInts, UserResponseIntegration{
							ID:     intMap["SmartCar"],
							Vendor: "SmartCar",
							Points: 1000,
						})
					}
				} else if services.ContainsString(ints, intMap["Tesla"]) {
					outInts = append(outInts, UserResponseIntegration{
						ID:     intMap["Tesla"],
						Vendor: "Tesla",
						Points: 4000,
					})
				} else if services.ContainsString(ints, intMap["SmartCar"]) {
					outInts = append(outInts, UserResponseIntegration{
						ID:     intMap["SmartCar"],
						Vendor: "SmartCar",
						Points: 1000,
					})
				}
			}
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

		tokens, err := models.TokenAllocations(
			models.TokenAllocationWhere.UserDeviceID.EQ(device.Id),
			qm.OrderBy(models.RewardColumns.IssuanceWeekID+" desc"),
		).All(c.Context(), r.DB().Reader)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve previously earned tokens.")
			return opaqueInternalError
		}
		tkns := big.NewInt(0)
		for _, t := range tokens {
			deviceTokens, bool := t.Tokens.Int64()
			if !bool {
				logger.Fatal().Msg("unable to convert weekly token allocation")
			}
			tkns.Add(tkns, big.NewInt(deviceTokens))
		}

		lvl := 1
		connectionStreak := 0
		disconnectionStreak := 0
		if len(rewards) > 0 {
			lvl = services.GetLevel(rewards[0].ConnectionStreak).Level
			connectionStreak = rewards[0].ConnectionStreak
			disconnectionStreak = rewards[0].DisconnectionStreak
		}

		outLi[i] = &UserResponseDevice{
			ID:                   device.Id,
			Points:               pts,
			Tokens:               tkns,
			ConnectedThisWeek:    activeThisWeek,
			IntegrationsThisWeek: outInts,
			LastActive:           maybeLastActive,
			ConnectionStreak:     connectionStreak,
			DisconnectionStreak:  disconnectionStreak,
			Level:                *getLevelResp(lvl),
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

func getLevelResp(level int) *UserResponseLevel {
	info := services.LevelInfos[level-1]
	var maxWeeks *int
	if level < 4 {
		maxWk := services.LevelInfos[level].MinWeeks - 1
		maxWeeks = &maxWk
	}
	return &UserResponseLevel{
		Number:       level,
		MinWeeks:     info.MinWeeks,
		MaxWeeks:     maxWeeks,
		StreakPoints: info.Points,
	}
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
	// Tokens is the total number of tokens that the device has earned across all weeks.
	Tokens *big.Int `json:"tokens" example:"5000"`
	// ConnectedThisWeek is true if we've seen activity from the device during the current issuance
	// week.
	ConnectedThisWeek bool `json:"connectedThisWeek" example:"true"`
	// IntegrationsThisWeek details the integrations we've seen active this week.
	IntegrationsThisWeek []UserResponseIntegration `json:"integrationsThisWeek"`
	// LastActive is the last time we saw activity from the vehicle.
	LastActive *time.Time `json:"lastActive,omitempty" example:"2022-04-12T09:23:01Z"`
	// ConnectionStreak is what we consider the streak of the device to be. This may not literally
	// be the number of consecutive connected weeks, because the user may disconnect for a week
	// without penalty, or have the connection streak reduced after three weeks of inactivity.
	ConnectionStreak int `json:"connectionStreak" example:"4"`
	// DisconnectionStreak is the number of consecutive issuance weeks that the device has been
	// disconnected. This number resets to 0 as soon as a device earns rewards for a certain week.
	DisconnectionStreak int `json:"disconnectionStreak,omitempty" example:"0"`
	// Level is the level 1-4 of the device. This is fully determined by ConnectionStreak.
	Level UserResponseLevel `json:"level"`
}

type UserResponseThisWeek struct {
	// Start is the timestamp of the start of the issuance week.
	Start time.Time `json:"start" example:"2022-04-11T05:00:00Z"`
	// End is the timestamp of the start of the next issuance week.
	End time.Time `json:"end" example:"2022-04-18T05:00:00Z"`
}

type UserResponseLevel struct {
	// Number is the level number 1-4
	Number int `json:"number" example:"2"`
	// MinWeeks is the minimum streak of weeks needed to enter this level.
	MinWeeks int `json:"minWeeks" example:"4"`
	// MaxWeeks is the last streak week at this level. In the next week, we enter the next level.
	MaxWeeks *int `json:"maxWeeks,omitempty" example:"20"`
	// StreakPoints is the number of points you earn per week at this level.
	StreakPoints int `json:"streakPoints" example:"1000"`
}

type UserResponseIntegration struct {
	// ID is the integration ID.
	ID string `json:"id" example:"27egBSLazAT7njT2VBjcISPIpiU"`
	// Vendor is the name of the integration vendor. At present, this uniquely determines the
	// integration.
	Vendor string `json:"vendor" example:"SmartCar"`
	// Points is the number of points a user earns for being connected with this integration
	// for a week.
	Points int `json:"points" example:"1000"`
}

// GetUserRewardsHistory godoc
// @Description  A summary of the user's rewards for past weeks.
// @Success      200 {object} controllers.HistoryResponse
// @Security     BearerAuth
// @Router       /user/history [get]
func (r *RewardsController) GetUserRewardsHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	devices, err := r.DevicesClient.ListUserDevicesForUser(c.Context(), &pb.ListUserDevicesForUserRequest{
		UserId: userID,
	})
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user's devices.")
		return opaqueInternalError
	}

	deviceIDs := make([]string, len(devices.UserDevices))
	for i := range devices.UserDevices {
		deviceIDs[i] = devices.UserDevices[i].Id
	}

	rs, err := models.Rewards(
		models.RewardWhere.UserID.EQ(userID),
		models.RewardWhere.UserDeviceID.IN(deviceIDs),
		qm.OrderBy(models.RewardColumns.IssuanceWeekID+" asc"),
	).All(c.Context(), r.DB().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving rewards.")
		return opaqueInternalError
	}

	if len(rs) == 0 {
		return c.JSON(HistoryResponse{})
	}

	minWeek := rs[0].IssuanceWeekID
	maxWeek := rs[len(rs)-1].IssuanceWeekID

	weeks := make([]HistoryResponseWeek, maxWeek-minWeek+1)
	for i := range weeks {
		weekNum := maxWeek - i
		weeks[i].Start = services.NumToWeekStart(weekNum)
		weeks[i].End = services.NumToWeekEnd(weekNum)
	}

	for _, s := range rs {
		distribution, err := models.TokenAllocations(models.TokenAllocationWhere.IssuanceWeekID.EQ(s.IssuanceWeekID), models.TokenAllocationWhere.UserDeviceID.EQ(s.UserDeviceID)).One(c.Context(), r.DB().Reader)
		if err != nil {
			logger.Err(err).Msg("unable to get device token allocation from table")
			return err
		}
		deviceTokens, bool := distribution.Tokens.Int64()
		if !bool {
			logger.Fatal().Msg("unable to convert weekly token allocation")
		}

		weeks[maxWeek-s.IssuanceWeekID].Points += s.StreakPoints + s.IntegrationPoints
		weeks[maxWeek-s.IssuanceWeekID].Tokens = big.NewInt(deviceTokens)
	}

	return c.JSON(HistoryResponse{Weeks: weeks})
}

type HistoryResponse struct {
	Weeks []HistoryResponseWeek `json:"weeks"`
}

type HistoryResponseWeek struct {
	// Start is the starting time of the issuance week.
	Start time.Time `json:"start" example:"2022-04-11T05:00:00Z"`
	// End is the starting time of the issuance week after this one.
	End time.Time `json:"end" example:"2022-04-18T05:00:00Z"`
	// Points is the number of points the user earned this week.
	Points int `json:"points" example:"4000"`
	// Tokens is the number of tokens the user earned this week.
	Tokens *big.Int `json:"tokens" example:"4000"`
}

type PointsDistributed struct {
	WeekStart time.Time `json:"week_start"`
	WeekEnd   time.Time `json:"week_end"`
	Points    int64     `json:"points,omitempty"`
	Tokens    int64     `json:"tokens,omitempty"`
}

// GetPointsThisWeek godoc
// @Description  Total number of points distributed to users this week
// @Success      200 {object} controllers.UserResponse
// @Security     BearerAuth
// @Router       /points [get]
func (r *RewardsController) GetPointsThisWeek(c *fiber.Ctx) error {
	now := time.Now()
	weekNum := services.GetWeekNum(now)

	pointsDistributed, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(weekNum)).One(c.Context(), r.DB().Reader)
	if err != nil {
		return err
	}

	return c.JSON(PointsDistributed{WeekStart: pointsDistributed.StartsAt, WeekEnd: pointsDistributed.EndsAt, Points: pointsDistributed.PointsDistributed.Int64})
}

// GetTokensThisWeek godoc
// @Description  Total number of tokens distributed to users this week
// @Success      200 {object} controllers.UserResponse
// @Security     BearerAuth
// @Router       /tokens [get]
func (r *RewardsController) GetTokensThisWeek(c *fiber.Ctx) error {
	now := time.Now()
	weekNum := services.GetWeekNum(now)

	pointsDistributed, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(weekNum)).One(c.Context(), r.DB().Reader)
	if err != nil {
		return err
	}
	tokensDistributed, bool := pointsDistributed.WeeklyTokenAllocation.Int64()
	if !bool {
		return err
	}

	return c.JSON(PointsDistributed{WeekStart: pointsDistributed.StartsAt, WeekEnd: pointsDistributed.EndsAt, Points: tokensDistributed})
}
