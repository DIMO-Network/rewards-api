package controllers

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/types/known/emptypb"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	pb_users "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/exp/slices"
)

type RewardsController struct {
	DB                db.Store
	Logger            *zerolog.Logger
	DataClient        services.DeviceDataClient
	DefinitionsClient pb_defs.DeviceDefinitionServiceClient
	DevicesClient     pb_devices.UserDeviceServiceClient
	UsersClient       pb_users.UserServiceClient
	Settings          *config.Settings
	Tokens            []*contracts.Token
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

	user, err := r.UsersClient.GetUser(c.Context(), &pb_users.GetUserRequest{Id: userID})
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User unknown.")
	}

	var addrBalance *big.Int

	if addr := user.EthereumAddress; addr != nil {
		addrBalance = big.NewInt(0)
		for _, tk := range r.Tokens {
			val, err := tk.BalanceOf(nil, common.HexToAddress(*addr))
			if err != nil {
				return err
			}
			addrBalance.Add(addrBalance, val)
		}
	}

	devicesReq := &pb_devices.ListUserDevicesForUserRequest{UserId: userID}
	if user.EthereumAddress != nil {
		devicesReq.EthereumAddress = *user.EthereumAddress
	}

	devices, err := r.DevicesClient.ListUserDevicesForUser(c.Context(), devicesReq)
	if err != nil {
		return err
	}

	allIntegrations, err := r.DefinitionsClient.GetIntegrations(c.Context(), &emptypb.Empty{})
	if err != nil {
		return opaqueInternalError
	}

	amMfrTokenToIntegration := make(map[uint64]*pb_defs.Integration)
	swIntegrsByTokenID := make(map[uint64]*pb_defs.Integration)

	for _, intDesc := range allIntegrations.Integrations {
		if intDesc.ManufacturerTokenId == 0 {
			// Must be a software integration. Sort after this loop.
			swIntegrsByTokenID[intDesc.TokenId] = intDesc
		} else {
			// Must be the integration associated with a manufacturer.
			amMfrTokenToIntegration[intDesc.ManufacturerTokenId] = intDesc
		}
	}

	outLi := make([]*UserResponseDevice, len(devices.UserDevices))
	userPts := 0
	userTokens := big.NewInt(0)

	for i, device := range devices.UserDevices {
		dlog := logger.With().Str("userDeviceId", device.Id).Logger()

		var maybeLastActive *time.Time
		lastActive, seen, err := r.DataClient.GetLastActivity(device.Id)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve last activity.")
			return opaqueInternalError
		}

		if seen {
			maybeLastActive = &lastActive
		}

		outInts := []UserResponseIntegration{}

		vehicleMinted := device.TokenId != nil

		vehicleIntegsWithSignals, err := r.DataClient.GetIntegrations(device.Id, weekStart, now)
		if err != nil {
			return opaqueInternalError
		}

		integSignalsThisWeek := utils.NewSet[string](vehicleIntegsWithSignals...)

		if ad := device.AftermarketDevice; ad != nil {
			// Want to see if this kind (right manufacturer) of device transmitted for this vehicle
			// this week.
			if ad.ManufacturerTokenId == 0 {
				return fmt.Errorf("aftermarket device %d does not have a manufacturer", ad.TokenId)
			}

			integr, ok := amMfrTokenToIntegration[ad.ManufacturerTokenId]
			if !ok {
				return fmt.Errorf("aftermarket device manufacturer %d does not have an associated integration", ad.ManufacturerTokenId)
			}

			uri := UserResponseIntegration{
				ID:                   integr.Id,
				Vendor:               integr.Vendor,
				DataThisWeek:         false,
				Points:               0,
				OnChainPairingStatus: "Paired",
			}

			if vehicleMinted && integSignalsThisWeek.Contains(integr.Id) {
				uri.Points = int(integr.Points)
				uri.DataThisWeek = true
			}

			outInts = append(outInts, uri)
		}

		if sd := device.SyntheticDevice; sd != nil {
			if sd.IntegrationTokenId == 0 {
				return fmt.Errorf("synthetic device %d does not have an integration", sd.IntegrationTokenId)
			}

			integr, ok := swIntegrsByTokenID[sd.IntegrationTokenId]
			if !ok {
				return fmt.Errorf("synthetic device %d has integration %d without metadata", sd.TokenId, sd.IntegrationTokenId)
			}

			uri := UserResponseIntegration{
				ID:                   integr.Id,
				Vendor:               integr.Vendor,
				DataThisWeek:         false,
				Points:               0,
				OnChainPairingStatus: "Unpaired",
			}

			if vehicleMinted && integSignalsThisWeek.Contains(integr.Id) {
				uri.Points = int(integr.Points)
				uri.DataThisWeek = true
				uri.OnChainPairingStatus = "Paired"
			}

			outInts = append(outInts, uri)
		}

		rewards, err := models.Rewards(
			models.RewardWhere.UserDeviceID.EQ(device.Id),
			qm.OrderBy(models.RewardColumns.IssuanceWeekID+" DESC"),
		).All(c.Context(), r.DB.DBS().Reader.DB)
		if err != nil {
			dlog.Err(err).Msg("Failed to retrieve previously earned rewards.")
			return opaqueInternalError
		}

		pts := 0
		for _, r := range rewards {
			pts += r.StreakPoints + r.AftermarketDevicePoints + r.SyntheticDevicePoints
		}

		userPts += pts

		tkns := big.NewInt(0)
		for _, t := range rewards {
			if t.AftermarketDeviceTokens.IsZero() && t.SyntheticDeviceTokens.IsZero() && t.StreakTokens.IsZero() {
				continue
			}
			tkns.Add(tkns, t.StreakTokens.Int(nil))
			tkns.Add(tkns, t.SyntheticDeviceTokens.Int(nil))
			tkns.Add(tkns, t.AftermarketDeviceTokens.Int(nil))
		}

		userTokens.Add(userTokens, tkns)

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
			ConnectedThisWeek:    slices.ContainsFunc(outInts, func(uri UserResponseIntegration) bool { return uri.Points > 0 }),
			IntegrationsThisWeek: outInts,
			LastActive:           maybeLastActive,
			ConnectionStreak:     connectionStreak,
			DisconnectionStreak:  disconnectionStreak,
			Level:                *getLevelResp(lvl),
			Minted:               device.TokenId != nil,
			OptedIn:              true,
		}
	}

	out := UserResponse{
		Points:        userPts,
		Tokens:        userTokens,
		WalletBalance: addrBalance,
		ThisWeek: UserResponseThisWeek{
			Start: weekStart,
			End:   services.NumToWeekEnd(weekNum),
		},
		Devices: outLi,
	}

	return c.JSON(out)
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
	// Tokens is the number of tokens the user has earned, across all devices and issuance
	// weeks.
	Tokens *big.Int `json:"tokens" example:"1105000000000000000000000" swaggertype:"number"`
	// WalletBalance is the number of tokens held in the users's wallet, if he has a wallet
	// attached to the present account.
	WalletBalance *big.Int `json:"walletBalance" example:"1105000000000000000000000" swaggertype:"number"`
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
	Tokens *big.Int `json:"tokens,omitempty" example:"5000" swaggertype:"number"`
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
	// Minted is true if the device has been minted on-chain.
	Minted bool `json:"minted"`
	// OptedIn is true if the user has agreed to the terms of service.
	OptedIn bool `json:"optedIn"`
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
	Vendor       string `json:"vendor" example:"SmartCar"`
	DataThisWeek bool   `json:"dataThisWeek"`
	// Points is the number of points a user earns for being connected with this integration
	// for a week.
	Points int `json:"points" example:"1000"`
	// OnChainPairingStatus is the on-chain pairing status of the integration.
	OnChainPairingStatus string `json:"onChainPairingStatus" enums:"Paired,Unpaired,NotApplicable" example:"Paired"`
}

// GetUserRewardsHistory godoc
// @Description  A summary of the user's rewards for past weeks.
// @Success      200 {object} controllers.HistoryResponse
// @Security     BearerAuth
// @Router       /user/history [get]
func (r *RewardsController) GetUserRewardsHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	user, err := r.UsersClient.GetUser(c.Context(), &pb_users.GetUserRequest{Id: userID})
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User unknown.")
	}

	devicesReq := &pb_devices.ListUserDevicesForUserRequest{UserId: userID}

	if user.EthereumAddress != nil {
		devicesReq.EthereumAddress = *user.EthereumAddress
	}

	devices, err := r.DevicesClient.ListUserDevicesForUser(c.Context(), devicesReq)
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user's devices.")
		return opaqueInternalError
	}

	deviceIDs := make([]string, len(devices.UserDevices))
	for i := range devices.UserDevices {
		deviceIDs[i] = devices.UserDevices[i].Id
	}

	rs, err := models.Rewards(
		models.RewardWhere.UserDeviceID.IN(deviceIDs),
		qm.OrderBy(models.RewardColumns.IssuanceWeekID+" asc"),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving rewards.")
		return opaqueInternalError
	}

	if len(rs) == 0 {
		return c.JSON(HistoryResponse{Weeks: []HistoryResponseWeek{}})
	}

	minWeek := rs[0].IssuanceWeekID
	maxWeek := rs[len(rs)-1].IssuanceWeekID

	weeks := make([]HistoryResponseWeek, maxWeek-minWeek+1)
	for i := range weeks {
		weekNum := maxWeek - i
		weeks[i].Start = services.NumToWeekStart(weekNum)
		weeks[i].End = services.NumToWeekEnd(weekNum)
		weeks[i].Tokens = big.NewInt(0)
	}

	tkns := big.NewInt(0)
	for _, r := range rs {
		weeks[maxWeek-r.IssuanceWeekID].Points += r.StreakPoints + r.AftermarketDevicePoints + r.SyntheticDevicePoints

		if r.AftermarketDeviceTokens.IsZero() && r.SyntheticDeviceTokens.IsZero() && r.StreakTokens.IsZero() {
			continue
		}

		tkns.Add(tkns, r.StreakTokens.Int(nil))
		tkns.Add(tkns, r.SyntheticDeviceTokens.Int(nil))
		tkns.Add(tkns, r.AftermarketDeviceTokens.Int(nil))

		weeks[maxWeek-r.IssuanceWeekID].Tokens.Add(weeks[maxWeek-r.IssuanceWeekID].Tokens, tkns)
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
	Tokens *big.Int `json:"tokens" example:"4000" swaggertype:"number"`
}
