package services

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/identity"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/DIMO-Network/shared/pkg/settings"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/mock/gomock"
)

const (
	macaron = "0x4c674ddE8189aEF6e3b58F5a36d7438b2b1f6Bc2"
	tesla   = "0xc4035Fecb1cc906130423EF05f9C20977F643722"
	autoPi  = "0x5e31bBc786D7bEd95216383787deA1ab0f1c1897"
	ruptela = "0xF26421509Efe92861a587482100c6d728aBf1CD0"
)

func MakeAddr(s string) common.Address {
	return crypto.PubkeyToAddress(crypto.ToECDSAUnsafe(crypto.Keccak256([]byte(s))).PublicKey)
}

const issuanceWeek = 30

type OldReward struct {
	ConnStreak        int
	DiscStreak        int
	UserDeviceTokenID int
}

type NewReward struct {
	TokenID                        int
	ConnStreak                     int
	DiscStreak                     int
	StreakPoints                   int
	RewardsReceiverEthereumAddress common.Address
	SyntheticDeviceID              int
	AftermarketDeviceID            int
	AftermarketDevicePoints        int
	SyntheticDevicePoints          int
}

type VIN struct {
	VIN        string
	FirstWeek  int
	FirstToken int
}

type Scenario struct {
	Name        string
	Previous    []OldReward
	Devices     []Device
	New         []NewReward
	PrevVIN     []VIN
	NewVIN      []VIN
	Description string
}

type RewardEvent struct {
	User      common.Address
	Value     *big.Int
	VehicleID *big.Int
}

var mkAddr = func(i int) common.Address {
	return common.BigToAddress(big.NewInt(int64(i)))
}

var mkVIN = func(i int) string {
	s := fmt.Sprintf("%017d", i)
	return s
}

func TestGetWeekNumForCron(t *testing.T) {
	ti, _ := time.Parse(time.RFC3339, "2022-02-07T05:00:02Z")
	if date.GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}

	ti, _ = time.Parse(time.RFC3339, "2022-02-07T04:58:44Z")
	if date.GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}
}

var teslaConn = common.HexToAddress("0xc4035Fecb1cc906130423EF05f9C20977F643722")

func TestStreak(t *testing.T) {
	ctx := context.Background()

	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	user1 := MakeAddr("user1")

	scens := []Scenario{
		{
			Name: "Level1TeslaGrow",
			Devices: []Device{
				{TokenID: 1, Owner: user1, VIN: mkVIN(1), DataSources: []string{tesla}, SDTokenID: 3, SDConnAddr: teslaConn},
			},
			Previous: []OldReward{
				{UserDeviceTokenID: 1, ConnStreak: 1, DiscStreak: 0},
			},
			New: []NewReward{
				{TokenID: 1, RewardsReceiverEthereumAddress: user1, ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 6000, SyntheticDeviceID: 3},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 1, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name:    "Level1TeslaDisconnected",
			Devices: []Device{},
			Previous: []OldReward{
				{ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{TokenID: 1, ConnStreak: 1, DiscStreak: 1, StreakPoints: 0, SyntheticDevicePoints: 0, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 1, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "AutoPiJoinLevel2",
			Devices: []Device{
				{TokenID: 1, Owner: user1, VIN: mkVIN(1), DataSources: []string{autoPi}, AMTokenID: 12, ManufacturerTokenID: 137},
			},
			Previous: []OldReward{
				{ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{TokenID: 1, ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, SyntheticDeviceID: 0, AftermarketDeviceID: 12, RewardsReceiverEthereumAddress: user1},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "LosingLevel3",
			Devices: []Device{
				{TokenID: 1, Owner: user1, VIN: mkVIN(1), DataSources: []string{}},
			},
			Previous: []OldReward{
				{ConnStreak: 22, DiscStreak: 2, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{TokenID: 1, ConnStreak: 4, DiscStreak: 3},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 1, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "BrandNewTesla",
			Devices: []Device{
				{TokenID: 1, Owner: user1, VIN: mkVIN(1), DataSources: []string{tesla}, SDTokenID: 3, SDConnAddr: teslaConn},
			},
			Previous: []OldReward{},
			New: []NewReward{
				{TokenID: 1, RewardsReceiverEthereumAddress: user1, ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 6000, SyntheticDeviceID: 3},
			},
			PrevVIN: []VIN{},
			NewVIN:  []VIN{{VIN: mkVIN(1), FirstToken: 1}},
		},
		{
			Name: "NewCopySameVIN",
			Devices: []Device{
				{TokenID: 1, VIN: mkVIN(1), DataSources: []string{}, Owner: MakeAddr("user1"), SDTokenID: 3, SDConnAddr: teslaConn},
				{TokenID: 2, VIN: mkVIN(1), DataSources: []string{tesla}, SDTokenID: 4, SDConnAddr: teslaConn, Owner: MakeAddr("user2")},
			},
			Previous: []OldReward{
				{ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{TokenID: 1, ConnStreak: 1, DiscStreak: 1},
				{TokenID: 2, RewardsReceiverEthereumAddress: MakeAddr("user2"), ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 6000, SyntheticDeviceID: 4},
			},
			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
			NewVIN:  []VIN{},
		},
		// {
		// 	Name: "Multiple HW Connections, paired on-chain is counted (autopi)",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{autoPiIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, SyntheticDeviceID: 0},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
		// {
		// 	Name: "Multiple HW Connections, paired on-chain is counted (macaron)",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{autoPiIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 2000, SyntheticDeviceID: 0},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
		// {
		// 	Name: "Combined HW and SW Connections, Macaron and Smartcar",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{smartcarIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142, SDTokenID: 1, SDIntegrationID: 3},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, AftermarketDevicePoints: 2000, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
		// {
		// 	Name: "Combined HW and SW Connections, AutoPi and Smartcar",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{smartcarIntegration, autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, SDTokenID: 1, SDIntegrationID: 3},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1, AftermarketDevicePoints: 6000},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
		// {
		// 	Name: "Combined HW and SW Connections, AutoPi and Smartcar-- AP does not transmit",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{smartcarIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, SDTokenID: 21, SDIntegrationID: 3},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 21},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
		// {
		// 	Name: "Multiple SW Connections, only one paired synthetic device",
		// 	Users: []User{
		// 		{ID: "User1", Address: mkAddr(1)},
		// 	},
		// 	Devices: []Device{
		// 		{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), DataSources: []string{smartcarIntegration, teslaIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142, SDTokenID: 1, SDIntegrationID: 3},
		// 	},
		// 	Previous: []OldReward{
		// 		{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
		// 	},
		// 	New: []NewReward{
		// 		{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
		// 	},
		// 	PrevVIN: []VIN{
		// 		{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
		// 	},
		// 	NewVIN: []VIN{},
		// },
	}

	for _, scen := range scens {
		t.Run(scen.Name, func(t *testing.T) {
			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

			_, err := models.Vins().DeleteAll(ctx, conn.DBS().Writer)
			assert.NoError(t, err)

			_, err = models.IssuanceWeeks().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

			lastWk := models.IssuanceWeek{ID: 0, JobStatus: models.IssuanceWeeksJobStatusFinished}
			err = lastWk.Insert(ctx, conn.DBS().Writer, boil.Infer())
			if err != nil {
				t.Fatal(err)
			}

			for _, lst := range scen.Previous {

				wk := models.IssuanceWeek{
					ID:        issuanceWeek - 1,
					JobStatus: models.IssuanceWeeksJobStatusFinished,
				}

				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
				assert.NoError(t, err)

				rw := models.Reward{
					IssuanceWeekID:      issuanceWeek - 1,
					ConnectionStreak:    lst.ConnStreak,
					DisconnectionStreak: lst.DiscStreak,
					UserDeviceTokenID:   lst.UserDeviceTokenID,
				}
				err = rw.Insert(ctx, conn.DBS().Writer, boil.Infer())
				if err != nil {
					t.Fatal(err)
				}
			}

			for _, v := range scen.PrevVIN {
				vw := models.Vin{
					Vin:                 v.VIN,
					FirstEarningWeek:    v.FirstWeek,
					FirstEarningTokenID: types.NewDecimal(decimal.New(int64(v.FirstToken), 0)),
				}

				err := vw.Insert(ctx, conn.DBS().Writer.DB, boil.Infer())
				assert.NoError(t, err)
			}

			transferService := NewTokenTransferService(&settings, nil, conn)

			ctrl := gomock.NewController(t)

			chClient := NewMockDeviceActivityClient(ctrl)
			devicesClient := NewMockDevicesClient(ctrl)
			identClient := NewMockIdentityClient(ctrl)

			var devices []*ch.Vehicle

			for _, v := range scen.Devices {
				if len(v.DataSources) != 0 {
					devices = append(devices, &ch.Vehicle{
						TokenID: int64(v.TokenID),
						Sources: v.DataSources,
					})

					m := map[string]any{
						"tokenId": v.TokenID,
						"owner":   v.Owner,
					}
					if v.SDTokenID != 0 {
						m["syntheticDevice"] = map[string]any{
							"tokenId": v.SDTokenID,
							"connection": map[string]any{
								"address": v.SDConnAddr,
							},
						}
					}
					if v.AMTokenID != 0 {
						m["aftermarketDevice"] = map[string]any{
							"tokenId":     v.AMTokenID,
							"beneficiary": v.Owner,
							"manufacturer": map[string]any{
								"tokenId": v.ManufacturerTokenID,
							},
						}

						devicesClient.EXPECT().GetVehicleByTokenIdFast(gomock.Any(), &pb_devices.GetVehicleByTokenIdFastRequest{
							TokenId: uint32(v.TokenID),
						}).Return(&pb_devices.GetVehicleByTokenIdFastResponse{
							Vin: v.VIN,
						}, nil)
					}

					var vd identity.VehicleDescription
					require.NoError(t, mapstructure.Decode(m, &vd))
					identClient.EXPECT().DescribeVehicle(uint64(v.TokenID)).Return(&vd, nil)
				}
			}

			chClient.EXPECT().DescribeActiveDevices(gomock.Any(), date.NumToWeekStart(issuanceWeek), date.NumToWeekEnd(issuanceWeek)).Return(devices, nil)

			rwBonusService := NewBaselineRewardService(&settings, transferService, chClient, identClient, issuanceWeek, &logger, nil)

			err = rwBonusService.assignPoints()
			if err != nil {
				t.Fatal(err)
			}

			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.UserDeviceTokenID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			var actual []NewReward

			for _, c := range rw {
				adID := 0
				if !c.AftermarketTokenID.IsZero() {
					a, _ := c.AftermarketTokenID.Int64()
					adID = int(a)
				}

				nr := NewReward{
					ConnStreak:                     c.ConnectionStreak,
					DiscStreak:                     c.DisconnectionStreak,
					StreakPoints:                   c.StreakPoints,
					AftermarketDevicePoints:        c.AftermarketDevicePoints,
					SyntheticDevicePoints:          c.SyntheticDevicePoints,
					SyntheticDeviceID:              c.SyntheticDeviceID.Int,
					TokenID:                        c.UserDeviceTokenID,
					AftermarketDeviceID:            adID,
					RewardsReceiverEthereumAddress: common.HexToAddress(c.RewardsReceiverEthereumAddress.String),
				}

				actual = append(actual, nr)
			}

			assert.ElementsMatch(t, scen.New, actual)

			vs, err := models.Vins(models.VinWhere.FirstEarningWeek.EQ(issuanceWeek)).All(ctx, conn.DBS().Reader.DB)
			require.NoError(t, err)

			actualv := []VIN{}

			for _, v := range vs {
				i, _ := v.FirstEarningTokenID.Int64()
				actualv = append(actualv, VIN{VIN: v.Vin, FirstToken: int(i)})
			}

			assert.ElementsMatch(t, scen.NewVIN, actualv)
		})
	}
}

// func TestBeneficiaryAddressSetForRewards(t *testing.T) {
// 	ctx := context.Background()

// 	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	logger := zerolog.Nop()

// 	cont, conn := utils.GetDbConnection(ctx, t, logger)
// 	defer testcontainers.CleanupContainer(t, cont)

// 	scens := []Scenario{
// 		{
// 			Name: "AutoPiJoinLevel1",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(2)},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN:      []VIN{},
// 			Description: "Should set beneficiary as rewards receiver",
// 		},
// 		{
// 			Name: "AutoPiJoinLevel2",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(1)},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN:      []VIN{},
// 			Description: "Should leave owner as rewards receiver if beneficiary is not set",
// 		},
// 		{
// 			Name: "AutoPiJoinLevel3",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(1).Bytes()},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(1)},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN:      []VIN{},
// 			Description: "Should leave reward receiver as owner if beneficiary address is same as owner",
// 		},
// 	}

// 	for _, scen := range scens {
// 		t.Run(scen.Name, func(t *testing.T) {
// 			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			_, err := models.Vins().DeleteAll(ctx, conn.DBS().Writer)
// 			assert.NoError(t, err)

// 			_, err = models.IssuanceWeeks().DeleteAll(ctx, conn.DBS().Writer)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			lastWk := models.IssuanceWeek{ID: 0, JobStatus: models.IssuanceWeeksJobStatusFinished}
// 			err = lastWk.Insert(ctx, conn.DBS().Writer, boil.Infer())
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			for _, lst := range scen.Previous {

// 				wk := models.IssuanceWeek{
// 					ID:        lst.Week,
// 					JobStatus: models.IssuanceWeeksJobStatusFinished,
// 				}

// 				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
// 				assert.NoError(t, err)

// 				rw := models.Reward{
// 					IssuanceWeekID:      lst.Week,
// 					ConnectionStreak:    lst.ConnStreak,
// 					DisconnectionStreak: lst.DiscStreak,
// 					UserDeviceTokenID:   lst.UserDeviceTokenID,
// 				}
// 				err = rw.Insert(ctx, conn.DBS().Writer, boil.Infer())
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}

// 			for _, v := range scen.PrevVIN {
// 				vw := models.Vin{
// 					Vin:                 v.VIN,
// 					FirstEarningWeek:    v.FirstWeek,
// 					FirstEarningTokenID: types.NewDecimal(decimal.New(int64(v.FirstToken), 0)),
// 				}

// 				err := vw.Insert(ctx, conn.DBS().Writer.DB, boil.Infer())
// 				assert.NoError(t, err)
// 			}

// 			transferService := NewTokenTransferService(&settings, nil, conn)

// 			ctrl := gomock.NewController(t)
// 			msc := NewMockIdentityClient(ctrl)
// 			msc.EXPECT().GetVehicleStakePoints(gomock.Any()).Return(0, nil).AnyTimes()

// 			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, msc, vinVCSrv, issuanceWeek, &logger)

// 			err = rwBonusService.assignPoints()
// 			assert.NoError(t, err)

// 			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceTokenID)).All(ctx, conn.DBS().Reader)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			var actual []NewReward

// 			for _, c := range rw {
// 				nr := NewReward{
// 					ConnStreak:              c.ConnectionStreak,
// 					DiscStreak:              c.DisconnectionStreak,
// 					StreakPoints:            c.StreakPoints,
// 					AftermarketDevicePoints: c.AftermarketDevicePoints,
// 					SyntheticDevicePoints:   c.SyntheticDevicePoints,
// 					TokenID:                 c.UserDeviceTokenID,
// 				}

// 				if c.UserEthereumAddress.Valid {
// 					nr.Address = common.HexToAddress(c.UserEthereumAddress.String)
// 				}

// 				if c.RewardsReceiverEthereumAddress.Valid {
// 					nr.RewardsReceiverEthereumAddress = common.HexToAddress(c.RewardsReceiverEthereumAddress.String)
// 				}

// 				actual = append(actual, nr)
// 			}
// 			assert.ElementsMatch(t, scen.New, actual, scen.Description)

// 			vs, err := models.Vins(models.VinWhere.FirstEarningWeek.EQ(issuanceWeek)).All(ctx, conn.DBS().Reader.DB)
// 			require.NoError(t, err)

// 			actualv := []VIN{}

// 			for _, v := range vs {
// 				i, _ := v.FirstEarningTokenID.Int64()
// 				actualv = append(actualv, VIN{VIN: v.Vin, FirstToken: int(i)})
// 			}

// 			assert.ElementsMatch(t, scen.NewVIN, actualv)
// 		})
// 	}
// }

// func TestBaselineIssuance(t *testing.T) {
// 	ctx := context.Background()

// 	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	settings.TransferBatchSize = 2

// 	logger := zerolog.Nop()

// 	cont, conn := utils.GetDbConnection(ctx, t, logger)
// 	defer testcontainers.CleanupContainer(t, cont)

// 	scens := []Scenario{
// 		{
// 			Name: "AutoPiJoinLevel1",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(2)},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN:      []VIN{},
// 			Description: "Should transfer to beneficiary set on the device",
// 		},
// 		{
// 			Name: "AutoPiJoinLevel2",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN:      []VIN{},
// 			Description: "Should transfer to owner when beneficiary is not set on the device",
// 		},
// 		{
// 			Name: "SmartCar1",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, SDTokenID: 1, SDIntegrationID: 3},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN: []VIN{},
// 		},
// 		{
// 			Name: "SmartCarAutopi1",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration, autoPiIntegration}, SDTokenID: 1, SDIntegrationID: 3, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1, AftermarketDevicePoints: 0},
// 			},
// 			PrevVIN: []VIN{
// 				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
// 			},
// 			NewVIN: []VIN{},
// 		},
// 		{
// 			Name: "NewCopySameVIN",
// 			Users: []User{
// 				{ID: "User2", Address: mkAddr(2)},
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(2), TokenID: 2, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{}},
// 				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, SDTokenID: 1, SDIntegrationID: 3},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New: []NewReward{
// 				{TokenID: 2, ConnStreak: 0, DiscStreak: 0, StreakPoints: 0, AftermarketDevicePoints: 0},
// 				{TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
// 			},
// 			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
// 			NewVIN:  []VIN{},
// 		},
// 		{
// 			Name: "EmptyVehicles",
// 			Users: []User{
// 				{ID: "User1", Address: mkAddr(1)},
// 			},
// 			Devices: []Device{
// 				{ID: mkID(2), TokenID: 1, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{}},
// 			},
// 			Previous: []OldReward{
// 				{Week: 4, ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
// 			},
// 			New:     []NewReward{},
// 			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
// 			NewVIN:  []VIN{},
// 		},
// 	}

// 	for _, scen := range scens {
// 		t.Run(scen.Name, func(t *testing.T) {
// 			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
// 			assert.NoError(t, err)

// 			_, err := models.Vins().DeleteAll(ctx, conn.DBS().Writer)
// 			assert.NoError(t, err)

// 			_, err = models.IssuanceWeeks().DeleteAll(ctx, conn.DBS().Writer)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			lastWk := models.IssuanceWeek{ID: 0, JobStatus: models.IssuanceWeeksJobStatusFinished}
// 			err = lastWk.Insert(ctx, conn.DBS().Writer, boil.Infer())
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			for _, lst := range scen.Previous {

// 				wk := models.IssuanceWeek{
// 					ID:        lst.Week,
// 					JobStatus: models.IssuanceWeeksJobStatusFinished,
// 				}

// 				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
// 				assert.NoError(t, err)

// 				rw := models.Reward{
// 					IssuanceWeekID:      lst.Week,
// 					ConnectionStreak:    lst.ConnStreak,
// 					DisconnectionStreak: lst.DiscStreak,
// 					UserDeviceTokenID:   lst.UserDeviceTokenID,
// 				}
// 				err = rw.Insert(ctx, conn.DBS().Writer, boil.Infer())
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}

// 			for _, v := range scen.PrevVIN {
// 				vw := models.Vin{
// 					Vin:                 v.VIN,
// 					FirstEarningWeek:    v.FirstWeek,
// 					FirstEarningTokenID: types.NewDecimal(decimal.New(int64(v.FirstToken), 0)),
// 				}

// 				err := vw.Insert(ctx, conn.DBS().Writer.DB, boil.Infer())
// 				assert.NoError(t, err)
// 			}

// 			config := mocks.NewTestConfig()
// 			config.Producer.Return.Successes = true
// 			config.Producer.Return.Errors = true
// 			producer := mocks.NewSyncProducer(t, config)

// 			var out []cloudevent.CloudEvent[transferData]

// 			checker := func(b2 []byte) error {
// 				var o cloudevent.CloudEvent[transferData]
// 				err := json.Unmarshal(b2, &o)
// 				require.NoError(t, err)
// 				out = append(out, o)
// 				return nil
// 			}

// 			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

// 			transferService := NewTokenTransferService(&settings, producer, conn)

// 			ctrl := gomock.NewController(t)
// 			msc := NewMockStakeChecker(ctrl)
// 			msc.EXPECT().GetVehicleStakePoints(gomock.Any()).AnyTimes().Return(0, nil)
// 			vinVCSrv := NewMockVINVCService(ctrl)
// 			vinVCSrv.EXPECT().GetConfirmedVINVCs(gomock.Any(), gomock.Any(), issuanceWeek).AnyTimes().Return(map[int64]struct{}{}, nil)

// 			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, msc, vinVCSrv, issuanceWeek, &logger)

// 			err = rwBonusService.BaselineIssuance()
// 			assert.NoError(t, err)

// 			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceTokenID)).All(ctx, conn.DBS().Reader)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			producer.Close()

// 			abi, err := contracts.RewardMetaData.GetAbi()
// 			require.NoError(t, err)

// 			rewards := []contracts.RewardTransferInfo{}

// 			for i := range out {
// 				b := out[i].Data.Data
// 				require.NoError(t, err)
// 				o, _ := abi.Methods["batchTransfer"].Inputs.Unpack(b[4:])
// 				rwrds := o[0].([]struct {
// 					User                       common.Address `json:"user"`
// 					VehicleId                  *big.Int       `json:"vehicleId"`           //nolint
// 					AftermarketDeviceId        *big.Int       `json:"aftermarketDeviceId"` //nolint
// 					ValueFromAftermarketDevice *big.Int       `json:"valueFromAftermarketDevice"`
// 					SyntheticDeviceId          *big.Int       `json:"syntheticDeviceId"` //nolint
// 					ValueFromSyntheticDevice   *big.Int       `json:"valueFromSyntheticDevice"`
// 					ConnectionStreak           *big.Int       `json:"connectionStreak"`
// 					ValueFromStreak            *big.Int       `json:"valueFromStreak"`
// 				})

// 				for _, r := range rwrds {
// 					reward := contracts.RewardTransferInfo{
// 						User:                       r.User,
// 						VehicleId:                  r.VehicleId,           //nolint
// 						AftermarketDeviceId:        r.AftermarketDeviceId, //nolint
// 						ValueFromAftermarketDevice: r.ValueFromAftermarketDevice,
// 						SyntheticDeviceId:          r.SyntheticDeviceId, //nolint
// 						ValueFromSyntheticDevice:   r.ValueFromSyntheticDevice,
// 						ConnectionStreak:           r.ConnectionStreak,
// 						ValueFromStreak:            r.ValueFromStreak,
// 					}
// 					if r.AftermarketDeviceId.Int64() == 0 || r.ValueFromAftermarketDevice.Int64() == 0 {
// 						reward.AftermarketDeviceId = &big.Int{}
// 						reward.ValueFromAftermarketDevice = &big.Int{}
// 					}
// 					if r.SyntheticDeviceId.Int64() == 0 || r.ValueFromSyntheticDevice.Int64() == 0 {
// 						reward.SyntheticDeviceId = &big.Int{}
// 						reward.ValueFromSyntheticDevice = &big.Int{}
// 					}

// 					if r.ConnectionStreak.Int64() == 0 {
// 						reward.ConnectionStreak = &big.Int{}
// 					}

// 					if r.ValueFromStreak.Int64() == 0 {
// 						reward.ValueFromStreak = &big.Int{}
// 					}
// 					rewards = append(rewards, reward)
// 				}
// 			}

// 			expected := []contracts.RewardTransferInfo{}
// 			user := common.HexToAddress(rw[0].RewardsReceiverEthereumAddress.String)

// 			if len(rw) > 0 && scen.Name != "EmptyVehicles" {
// 				expected = append(expected, contracts.RewardTransferInfo{
// 					User:                       user,
// 					VehicleId:                  big.NewInt(int64(rw[0].UserDeviceTokenID)),
// 					AftermarketDeviceId:        utils.NullDecimalToIntDefaultZero(rw[0].AftermarketTokenID),
// 					ValueFromAftermarketDevice: utils.NullDecimalToIntDefaultZero(rw[0].AftermarketDeviceTokens),
// 					SyntheticDeviceId:          big.NewInt(int64(rw[0].SyntheticDeviceID.Int)),
// 					ValueFromSyntheticDevice:   utils.NullDecimalToIntDefaultZero(rw[0].SyntheticDeviceTokens),
// 					ConnectionStreak:           big.NewInt(int64(rw[0].ConnectionStreak)),
// 					ValueFromStreak:            utils.NullDecimalToIntDefaultZero(rw[0].StreakTokens),
// 				})
// 			}
// 			assert.ElementsMatch(t, expected, rewards)
// 		})
// 	}
// }

type Device struct {
	TokenID             int
	Owner               common.Address
	VIN                 string
	DataSources         []string
	AMTokenID           int
	ManufacturerTokenID int
	Beneficiary         []byte
	SDTokenID           uint64
	SDConnAddr          common.Address
}
