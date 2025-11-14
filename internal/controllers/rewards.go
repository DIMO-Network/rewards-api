package controllers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	att_types "github.com/DIMO-Network/attestation-api/pkg/types"
	"github.com/DIMO-Network/cloudevent"
	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	pb_fetch "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/constants"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/identity"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type RewardsController struct {
	DB            db.Store
	Logger        *zerolog.Logger
	ChClient      *ch.Client
	DevicesClient pb_devices.UserDeviceServiceClient
	IdentClient   *identity.Client
	Settings      *config.Settings
	FetchClient   pb_fetch.FetchServiceClient
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}

var zeroAddr common.Address

const ethAddrClaimName = "ethereum_address"

func GetTokenEthAddr(c *fiber.Ctx) (common.Address, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims) // These should never fail.
	ethAddrAny, ok := claims[ethAddrClaimName]
	if !ok {
		return zeroAddr, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("No %s claim in JWT.", ethAddrClaimName))
	}
	ethAddrStr, ok := ethAddrAny.(string) // These might
	if !ok {
		return zeroAddr, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Claim %s has unexpected type %T.", ethAddrClaimName, ethAddrAny))
	}
	if !common.IsHexAddress(ethAddrStr) {
		return zeroAddr, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Claim %s is not a valid Ethereum address.", ethAddrClaimName))
	}
	return common.HexToAddress(ethAddrStr), nil
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
	weekNum := date.GetWeekNum(now)
	weekStart := date.NumToWeekStart(weekNum)

	userAddr, err := GetTokenEthAddr(c)
	if err != nil {
		return err
	}

	addrBalance := big.NewInt(0)

	vehicleDescrs, err := r.IdentClient.GetVehicles(userAddr)
	if err != nil {
		return err
	}

	var vehicleIDs []uint64
	for _, ud := range vehicleDescrs {
		vehicleIDs = append(vehicleIDs, uint64(ud.TokenID))
	}

	sourcesByTokenID := make(map[uint64][]string)

	vehicleSources, err := r.ChClient.GetSourcesForVehicles(c.Context(), vehicleIDs, weekStart, now)
	if err != nil {
		return err
	}

	for _, v := range vehicleSources {
		sourcesByTokenID[uint64(v.TokenID)] = v.Sources
	}

	outLi := make([]*UserResponseDevice, 0, len(vehicleDescrs))
	userPts := 0
	userTokens := big.NewInt(0)

	for _, device := range vehicleDescrs {
		dlog := logger.With().Int("vehicleId", device.TokenID).Logger()

		outInts := []*UserResponseIntegration{}

		if ad := device.AftermarketDevice; ad != nil {
			conn, ok := constants.ConnsByMfrId[ad.Manufacturer.TokenID]
			if ok {
				uri := UserResponseIntegration{
					ID:                   conn.LegacyID,
					Vendor:               conn.LegacyVendor,
					DataThisWeek:         false,
					Points:               0,
					OnChainPairingStatus: "Paired",
				}

				if slices.Contains(sourcesByTokenID[uint64(device.TokenID)], conn.Address.Hex()) {
					uri.DataThisWeek = true
					uri.Points = int(conn.Points)
				}

				outInts = append(outInts, &uri)
			}
		}

		if sd := device.SyntheticDevice; sd != nil {
			conn, ok := constants.ConnsByAddr[sd.Connection.Address]
			if ok {
				uri := UserResponseIntegration{
					ID:                   conn.LegacyID,
					Vendor:               conn.LegacyVendor,
					DataThisWeek:         false,
					Points:               0,
					OnChainPairingStatus: "Paired",
				}

				if slices.Contains(sourcesByTokenID[uint64(device.TokenID)], conn.Address.Hex()) {
					uri.DataThisWeek = true
					uri.Points = int(conn.Points)
				}

				outInts = append(outInts, &uri)
			}
		}

		vinConfirmed := false

		// One last attempt, try the VIN VC.
		ce, err := r.FetchClient.GetLatestCloudEvent(c.Context(), &pb_fetch.GetLatestCloudEventRequest{
			Options: &pb_fetch.SearchOptions{
				Type:        &wrapperspb.StringValue{Value: cloudevent.TypeAttestation},
				DataVersion: &wrapperspb.StringValue{Value: "vin/v1.0"},
				Subject:     &wrapperspb.StringValue{Value: cloudevent.ERC721DID{ChainID: 137, ContractAddress: common.HexToAddress("0xbA5738a18d83D41847dfFbDC6101d37C69c9B0cF"), TokenID: big.NewInt(int64(device.TokenID))}.String()},
				Source:      &wrapperspb.StringValue{Value: common.HexToAddress("0x49eAf63eD94FEf3d40692862Eee2C8dB416B1a5f").Hex()},
			},
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				vtf, err := r.DevicesClient.GetVehicleByTokenIdFast(c.Context(), &pb_devices.GetVehicleByTokenIdFastRequest{
					TokenId: uint32(device.TokenID),
				})
				if err != nil {
					if status.Code(err) != codes.NotFound {
						// Some intermittent error.
						return fmt.Errorf("failed to grab vehicle %d: %w", device.TokenID, err)
					}
					// Otherwise, just leave vinConfirmed as false.
				} else if vtf.Vin != "" {
					vinConfirmed = true
				}
			} else {
				return fmt.Errorf("failed to retrieve VIN attestation for vehicle %d: %w", device.TokenID, err)
			}
		} else {
			var cred att_types.Credential
			if err := json.Unmarshal(ce.CloudEvent.Data, &cred); err == nil {
				var vs att_types.VINSubject
				if err := json.Unmarshal(cred.CredentialSubject, &vs); err == nil {
					vinConfirmed = true
				}
			}
		}

		if !vinConfirmed {
			for _, oi := range outInts {
				oi.DataThisWeek = false
				oi.Points = 0
			}
		}

		rewards, err := models.Rewards(
			models.RewardWhere.UserDeviceTokenID.EQ(device.TokenID),
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

		tokTemp := uint64(device.TokenID)

		outLi = append(outLi, &UserResponseDevice{
			TokenID:              &tokTemp,
			Points:               pts,
			Tokens:               tkns,
			ConnectedThisWeek:    slices.ContainsFunc(outInts, func(uri *UserResponseIntegration) bool { return uri.Points > 0 }),
			IntegrationsThisWeek: outInts,
			LastActive:           nil, // Unused, we think.
			ConnectionStreak:     connectionStreak,
			DisconnectionStreak:  disconnectionStreak,
			Level:                *getLevelResp(lvl),
			Minted:               true,
			OptedIn:              true,
			VINConfirmed:         vinConfirmed,
		})
	}

	out := UserResponse{
		Points:        userPts,
		Tokens:        userTokens,
		WalletBalance: addrBalance,
		ThisWeek: UserResponseThisWeek{
			Start: weekStart,
			End:   date.NumToWeekEnd(weekNum),
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
	// TokenID is the NFT token id for minted vehicles.
	TokenID *uint64 `json:"tokenId,omitempty" example:"37"`
	// Points is the total number of points that the device has earned across all weeks.
	Points int `json:"points" example:"5000"`
	// Tokens is the total number of tokens that the device has earned across all weeks.
	Tokens *big.Int `json:"tokens,omitempty" example:"5000" swaggertype:"number"`
	// ConnectedThisWeek is true if we've seen activity from the device during the current issuance
	// week.
	ConnectedThisWeek bool `json:"connectedThisWeek" example:"true"`
	// IntegrationsThisWeek details the integrations we've seen active this week.
	IntegrationsThisWeek []*UserResponseIntegration `json:"integrationsThisWeek"`
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
	OptedIn      bool `json:"optedIn"`
	VINConfirmed bool `json:"vinConfirmed"`
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

	userAddr, err := GetTokenEthAddr(c)
	if err != nil {
		return err
	}

	devicesReq := &pb_devices.ListUserDevicesForUserRequest{
		UserId:          userID,
		EthereumAddress: userAddr.Hex(),
	}

	devices, err := r.DevicesClient.ListUserDevicesForUser(c.Context(), devicesReq)
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve user's devices.")
		return opaqueInternalError
	}

	vehicleTokenIDs := make([]int, 0, len(devices.UserDevices))
	for _, ud := range devices.UserDevices {
		if ud.TokenId != nil {
			vehicleTokenIDs = append(vehicleTokenIDs, int(*ud.TokenId))
		}
	}

	rs, err := models.Rewards(
		models.RewardWhere.UserDeviceTokenID.IN(vehicleTokenIDs),
		qm.OrderBy(models.RewardColumns.IssuanceWeekID+" ASC"),
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
		weeks[i].Start = date.NumToWeekStart(weekNum)
		weeks[i].End = date.NumToWeekEnd(weekNum)
		weeks[i].Tokens = big.NewInt(0)
	}

	for _, r := range rs {
		weeks[maxWeek-r.IssuanceWeekID].Points += r.StreakPoints + r.AftermarketDevicePoints + r.SyntheticDevicePoints

		if r.AftermarketDeviceTokens.IsZero() && r.SyntheticDeviceTokens.IsZero() && r.StreakTokens.IsZero() {
			continue
		}

		tkns := big.NewInt(0)
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
