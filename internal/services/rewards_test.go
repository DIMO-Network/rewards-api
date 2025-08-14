package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/DIMO-Network/cloudevent"
	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/DIMO-Network/shared/pkg/settings"
	"github.com/IBM/sarama/mocks"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const issuanceWeek = 5

type User struct {
	ID       string
	Address  common.Address
	Code     string
	CodeUsed string
}

type OldReward struct {
	Week              int
	DeviceID          string
	UserID            string
	ConnStreak        int
	DiscStreak        int
	UserDeviceTokenID int64
}

type NewReward struct {
	DeviceID                       string
	TokenID                        int
	Address                        common.Address
	ConnStreak                     int
	DiscStreak                     int
	StreakPoints                   int
	RewardsReceiverEthereumAddress common.Address
	SyntheticDeviceID              int
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
	Users       []User
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

var mkID = func(i int) string {
	return fmt.Sprintf("% 27d", i)
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

func TestStreak(t *testing.T) {
	ctx := context.Background()

	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	scens := []Scenario{
		{
			Name: "Level1SmartcarGrow",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Level1SmartcarDisconnected",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 1, StreakPoints: 0, SyntheticDevicePoints: 0, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "AutoPiJoinLevel2",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "LosingLevel3",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 22, DiscStreak: 2, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 4, DiscStreak: 3, StreakPoints: 0, SyntheticDevicePoints: 0, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "BrandNewTesla",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{teslaIntegration}, SDTokenID: 3, SDIntegrationID: 2},
			},
			Previous: []OldReward{},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 4000, SyntheticDeviceID: 3},
			},
			PrevVIN: []VIN{},
			NewVIN:  []VIN{{VIN: mkVIN(1), FirstToken: 1}},
		},
		{
			Name: "NewCopySameVIN",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
				{ID: "User2", Address: mkAddr(2)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}},
				{ID: mkID(2), TokenID: 2, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{teslaIntegration}, SDTokenID: 2, SDIntegrationID: 2},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 1, StreakPoints: 0, AftermarketDevicePoints: 0},
				{DeviceID: mkID(2), TokenID: 2, Address: mkAddr(2), ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 4000, SyntheticDeviceID: 2},
			},
			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
			NewVIN:  []VIN{},
		},
		{
			Name: "Multiple HW Connections, paired on-chain is counted (autopi)",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Multiple HW Connections, paired on-chain is counted (macaron)",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 2000, SyntheticDeviceID: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Combined HW and SW Connections, Macaron and Smartcar",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration, macaronIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, AftermarketDevicePoints: 2000, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Combined HW and SW Connections, AutoPi and Smartcar",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration, autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1, AftermarketDevicePoints: 6000},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Combined HW and SW Connections, AutoPi and Smartcar-- AP does not transmit",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, SDTokenID: 21, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 21},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "Multiple SW Connections, only one paired synthetic device",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration, teslaIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 142, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
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
					ID:        lst.Week,
					JobStatus: models.IssuanceWeeksJobStatusFinished,
				}

				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
				assert.NoError(t, err)

				rw := models.Reward{
					IssuanceWeekID:      lst.Week,
					UserDeviceID:        lst.DeviceID,
					ConnectionStreak:    lst.ConnStreak,
					DisconnectionStreak: lst.DiscStreak,
					UserDeviceTokenID:   types.NewNullDecimal(decimal.New(lst.UserDeviceTokenID, 0)),
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
			msc := NewMockStakeChecker(ctrl)
			msc.EXPECT().GetVehicleStakePoints(gomock.Any()).AnyTimes().Return(0, nil)
			vinVCSrv := NewMockVINVCService(ctrl)
			vinVCSrv.EXPECT().GetConfirmedVINVCs(gomock.Any(), gomock.Any(), issuanceWeek).AnyTimes().Return(map[int64]struct{}{}, nil)

			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, msc, vinVCSrv, issuanceWeek, &logger)

			err = rwBonusService.assignPoints()
			if err != nil {
				t.Fatal(err)
			}

			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			var actual []NewReward

			for _, c := range rw {
				nr := NewReward{
					DeviceID:                c.UserDeviceID,
					ConnStreak:              c.ConnectionStreak,
					DiscStreak:              c.DisconnectionStreak,
					StreakPoints:            c.StreakPoints,
					AftermarketDevicePoints: c.AftermarketDevicePoints,
					SyntheticDevicePoints:   c.SyntheticDevicePoints,
					SyntheticDeviceID:       c.SyntheticDeviceID.Int,
				}

				if !c.UserDeviceTokenID.IsZero() {
					n, _ := c.UserDeviceTokenID.Int64()
					nr.TokenID = int(n)
				}

				if c.UserEthereumAddress.Valid {
					nr.Address = common.HexToAddress(c.UserEthereumAddress.String)
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

func TestBeneficiaryAddressSetForRewards(t *testing.T) {
	ctx := context.Background()

	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	scens := []Scenario{
		{
			Name: "AutoPiJoinLevel1",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(2)},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN:      []VIN{},
			Description: "Should set beneficiary as rewards receiver",
		},
		{
			Name: "AutoPiJoinLevel2",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(1)},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN:      []VIN{},
			Description: "Should leave owner as rewards receiver if beneficiary is not set",
		},
		{
			Name: "AutoPiJoinLevel3",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(1).Bytes()},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(1)},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN:      []VIN{},
			Description: "Should leave reward receiver as owner if beneficiary address is same as owner",
		},
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
					ID:        lst.Week,
					JobStatus: models.IssuanceWeeksJobStatusFinished,
				}

				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
				assert.NoError(t, err)

				rw := models.Reward{
					IssuanceWeekID:      lst.Week,
					UserDeviceID:        lst.DeviceID,
					ConnectionStreak:    lst.ConnStreak,
					DisconnectionStreak: lst.DiscStreak,
					UserDeviceTokenID:   types.NewNullDecimal(decimal.New(lst.UserDeviceTokenID, 0)),
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
			msc := NewMockStakeChecker(ctrl)
			msc.EXPECT().GetVehicleStakePoints(gomock.Any()).Return(0, nil).AnyTimes()
			vinVCSrv := NewMockVINVCService(ctrl)
			vinVCSrv.EXPECT().GetConfirmedVINVCs(gomock.Any(), gomock.Any(), issuanceWeek).AnyTimes().Return(map[int64]struct{}{}, nil)

			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, msc, vinVCSrv, issuanceWeek, &logger)

			err = rwBonusService.assignPoints()
			assert.NoError(t, err)

			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			var actual []NewReward

			for _, c := range rw {
				nr := NewReward{
					DeviceID:                c.UserDeviceID,
					ConnStreak:              c.ConnectionStreak,
					DiscStreak:              c.DisconnectionStreak,
					StreakPoints:            c.StreakPoints,
					AftermarketDevicePoints: c.AftermarketDevicePoints,
					SyntheticDevicePoints:   c.SyntheticDevicePoints,
				}

				if !c.UserDeviceTokenID.IsZero() {
					n, _ := c.UserDeviceTokenID.Int64()
					nr.TokenID = int(n)
				}

				if c.UserEthereumAddress.Valid {
					nr.Address = common.HexToAddress(c.UserEthereumAddress.String)
				}

				if c.RewardsReceiverEthereumAddress.Valid {
					nr.RewardsReceiverEthereumAddress = common.HexToAddress(c.RewardsReceiverEthereumAddress.String)
				}

				actual = append(actual, nr)
			}
			assert.ElementsMatch(t, scen.New, actual, scen.Description)

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

func TestBaselineIssuance(t *testing.T) {
	ctx := context.Background()

	settings, err := settings.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}
	settings.TransferBatchSize = 2

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	scens := []Scenario{
		{
			Name: "AutoPiJoinLevel1",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000, RewardsReceiverEthereumAddress: mkAddr(2)},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN:      []VIN{},
			Description: "Should transfer to beneficiary set on the device",
		},
		{
			Name: "AutoPiJoinLevel2",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, AftermarketDevicePoints: 6000},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN:      []VIN{},
			Description: "Should transfer to owner when beneficiary is not set on the device",
		},
		{
			Name: "SmartCar1",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "SmartCarAutopi1",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration, autoPiIntegration}, SDTokenID: 1, SDIntegrationID: 3, AMTokenID: 12, AMSerial: ksuid.New().String(), ManufacturerTokenID: 137, Beneficiary: mkAddr(2).Bytes()},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1, AftermarketDevicePoints: 0},
			},
			PrevVIN: []VIN{
				{VIN: mkVIN(1), FirstWeek: 4, FirstToken: 1},
			},
			NewVIN: []VIN{},
		},
		{
			Name: "NewCopySameVIN",
			Users: []User{
				{ID: "User2", Address: mkAddr(2)},
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(2), TokenID: 2, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{}},
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, SDTokenID: 1, SDIntegrationID: 3},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New: []NewReward{
				{DeviceID: mkID(2), ConnStreak: 0, DiscStreak: 0, StreakPoints: 0, AftermarketDevicePoints: 0, TokenID: 2},
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 0, SyntheticDevicePoints: 1000, SyntheticDeviceID: 1},
			},
			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
			NewVIN:  []VIN{},
		},
		{
			Name: "EmptyVehicles",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(2), TokenID: 1, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{}},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 0, UserDeviceTokenID: 1},
			},
			New:     []NewReward{},
			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
			NewVIN:  []VIN{},
		},
	}

	for _, scen := range scens {
		t.Run(scen.Name, func(t *testing.T) {
			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
			assert.NoError(t, err)

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
					ID:        lst.Week,
					JobStatus: models.IssuanceWeeksJobStatusFinished,
				}

				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
				assert.NoError(t, err)

				rw := models.Reward{
					IssuanceWeekID:      lst.Week,
					UserDeviceID:        lst.DeviceID,
					ConnectionStreak:    lst.ConnStreak,
					DisconnectionStreak: lst.DiscStreak,
					UserDeviceTokenID:   types.NewNullDecimal(decimal.New(lst.UserDeviceTokenID, 0)),
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

			config := mocks.NewTestConfig()
			config.Producer.Return.Successes = true
			config.Producer.Return.Errors = true
			producer := mocks.NewSyncProducer(t, config)

			var out []cloudevent.CloudEvent[transferData]

			checker := func(b2 []byte) error {
				var o cloudevent.CloudEvent[transferData]
				err := json.Unmarshal(b2, &o)
				require.NoError(t, err)
				out = append(out, o)
				return nil
			}

			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

			transferService := NewTokenTransferService(&settings, producer, conn)

			ctrl := gomock.NewController(t)
			msc := NewMockStakeChecker(ctrl)
			msc.EXPECT().GetVehicleStakePoints(gomock.Any()).AnyTimes().Return(0, nil)
			vinVCSrv := NewMockVINVCService(ctrl)
			vinVCSrv.EXPECT().GetConfirmedVINVCs(gomock.Any(), gomock.Any(), issuanceWeek).AnyTimes().Return(map[int64]struct{}{}, nil)

			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, msc, vinVCSrv, issuanceWeek, &logger)

			err = rwBonusService.BaselineIssuance()
			assert.NoError(t, err)

			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			producer.Close()

			abi, err := contracts.RewardMetaData.GetAbi()
			require.NoError(t, err)

			rewards := []contracts.RewardTransferInfo{}

			for i := range out {
				b := out[i].Data.Data
				require.NoError(t, err)
				o, _ := abi.Methods["batchTransfer"].Inputs.Unpack(b[4:])
				rwrds := o[0].([]struct {
					User                       common.Address `json:"user"`
					VehicleId                  *big.Int       `json:"vehicleId"`           //nolint
					AftermarketDeviceId        *big.Int       `json:"aftermarketDeviceId"` //nolint
					ValueFromAftermarketDevice *big.Int       `json:"valueFromAftermarketDevice"`
					SyntheticDeviceId          *big.Int       `json:"syntheticDeviceId"` //nolint
					ValueFromSyntheticDevice   *big.Int       `json:"valueFromSyntheticDevice"`
					ConnectionStreak           *big.Int       `json:"connectionStreak"`
					ValueFromStreak            *big.Int       `json:"valueFromStreak"`
				})

				for _, r := range rwrds {
					reward := contracts.RewardTransferInfo{
						User:                       r.User,
						VehicleId:                  r.VehicleId,           //nolint
						AftermarketDeviceId:        r.AftermarketDeviceId, //nolint
						ValueFromAftermarketDevice: r.ValueFromAftermarketDevice,
						SyntheticDeviceId:          r.SyntheticDeviceId, //nolint
						ValueFromSyntheticDevice:   r.ValueFromSyntheticDevice,
						ConnectionStreak:           r.ConnectionStreak,
						ValueFromStreak:            r.ValueFromStreak,
					}
					if r.AftermarketDeviceId.Int64() == 0 || r.ValueFromAftermarketDevice.Int64() == 0 {
						reward.AftermarketDeviceId = &big.Int{}
						reward.ValueFromAftermarketDevice = &big.Int{}
					}
					if r.SyntheticDeviceId.Int64() == 0 || r.ValueFromSyntheticDevice.Int64() == 0 {
						reward.SyntheticDeviceId = &big.Int{}
						reward.ValueFromSyntheticDevice = &big.Int{}
					}

					if r.ConnectionStreak.Int64() == 0 {
						reward.ConnectionStreak = &big.Int{}
					}

					if r.ValueFromStreak.Int64() == 0 {
						reward.ValueFromStreak = &big.Int{}
					}
					rewards = append(rewards, reward)
				}
			}

			expected := []contracts.RewardTransferInfo{}
			user := common.HexToAddress(rw[0].RewardsReceiverEthereumAddress.String)

			if len(rw) > 0 && scen.Name != "EmptyVehicles" {
				expected = append(expected, contracts.RewardTransferInfo{
					User:                       user,
					VehicleId:                  rw[0].UserDeviceTokenID.Int(nil),
					AftermarketDeviceId:        utils.NullDecimalToIntDefaultZero(rw[0].AftermarketTokenID),
					ValueFromAftermarketDevice: utils.NullDecimalToIntDefaultZero(rw[0].AftermarketDeviceTokens),
					SyntheticDeviceId:          big.NewInt(int64(rw[0].SyntheticDeviceID.Int)),
					ValueFromSyntheticDevice:   utils.NullDecimalToIntDefaultZero(rw[0].SyntheticDeviceTokens),
					ConnectionStreak:           big.NewInt(int64(rw[0].ConnectionStreak)),
					ValueFromStreak:            utils.NullDecimalToIntDefaultZero(rw[0].StreakTokens),
				})
			}
			assert.ElementsMatch(t, expected, rewards)
		})
	}
}

type Device struct {
	ID                  string
	TokenID             int
	UserID              string
	VIN                 string
	IntsWithData        []string
	AMTokenID           int
	AMSerial            string
	ManufacturerTokenID int
	Beneficiary         []byte
	SDTokenID           uint64
	SDIntegrationID     uint64
}

type Views struct {
	devices []Device
}

const (
	autoPiIntegration   = "2LFD6DXuGRdVucJO1a779kEUiYi"
	teslaIntegration    = "2LFQOgsYd5MEmRNBnsYXKp0QHC3"
	smartcarIntegration = "2LFSA81Oo4agy0y4NvP7f6hTdgs"
	macaronIntegration  = "2ULfuC8U9dOqRshZBAi0lMM1Rrx"
)

func (v Views) DescribeActiveDevices(ctx context.Context, _, _ time.Time) ([]*ch.Vehicle, error) {
	out := []*ch.Vehicle{}
	for _, d := range v.devices {
		if len(d.IntsWithData) == 0 {
			continue
		}
		out = append(out, &ch.Vehicle{
			TokenID: int64(d.TokenID), Integrations: d.IntsWithData,
		})
	}
	return out, nil
}

type FakeDefClient struct {
}

func (d *FakeDefClient) GetIntegrations(_ context.Context, _ *emptypb.Empty, _ ...grpc.CallOption) (*pb_defs.GetIntegrationResponse, error) {
	return &pb_defs.GetIntegrationResponse{Integrations: []*pb_defs.Integration{
		{Id: autoPiIntegration, Vendor: "AutoPi", ManufacturerTokenId: 137, Points: 6000},
		{Id: teslaIntegration, Vendor: "Tesla", Points: 4000, ManufacturerTokenId: 0, TokenId: 2},
		{Id: smartcarIntegration, Vendor: "SmartCar", Points: 1000, ManufacturerTokenId: 0, TokenId: 3},
		{Id: macaronIntegration, Vendor: "Macaron", ManufacturerTokenId: 142, Points: 2000},
	}}, nil
}

type FakeDevClient struct {
	users   []User
	devices []Device
}

var zeroAddr common.Address

func (d *FakeDevClient) GetUserDeviceByTokenId(_ context.Context, in *pb_devices.GetUserDeviceByTokenIdRequest, _ ...grpc.CallOption) (*pb_devices.UserDevice, error) {
	for _, ud := range d.devices {
		if int64(ud.TokenID) != in.TokenId {
			continue
		}

		var tk *uint64
		if ud.TokenID != 0 {
			t := uint64(ud.TokenID)
			tk = &t
		}

		var owner []byte
		for _, u := range d.users {
			if u.ID == ud.UserID {
				if u.Address != zeroAddr {
					owner = u.Address.Bytes()
				}
				break
			}
		}

		var vin *string
		if ud.VIN != "" {
			vin = &ud.VIN
		}

		ud2 := &pb_devices.UserDevice{
			Id:           ud.ID,
			TokenId:      tk,
			Vin:          vin,
			VinConfirmed: vin != nil,
			OwnerAddress: owner,
		}

		for _, i := range ud.IntsWithData {
			ud2.Integrations = append(ud2.Integrations, &pb_devices.UserDeviceIntegration{
				Id: i,
			})
		}

		if ud.AMTokenID != 0 {
			ud2.AftermarketDevice = &pb_devices.AftermarketDevice{
				Serial:              ud.AMSerial,
				UserId:              &ud.ID,
				OwnerAddress:        owner,
				TokenId:             *tk,
				ManufacturerTokenId: uint64(ud.ManufacturerTokenID),
				Beneficiary:         owner,
			}
			if len(ud.Beneficiary) != 0 {
				ud2.AftermarketDevice.Beneficiary = ud.Beneficiary
				ud2.AftermarketDeviceBeneficiaryAddress = ud.Beneficiary //nolint:staticcheck
			}

			ud2.AftermarketDeviceTokenId = &ud2.AftermarketDevice.TokenId //nolint:staticcheck
		}

		if ud.SDTokenID != 0 && ud.SDIntegrationID != 0 {
			ud2.SyntheticDevice = &pb_devices.SyntheticDevice{
				TokenId:            ud.SDTokenID,
				IntegrationTokenId: ud.SDIntegrationID,
			}
		}

		return ud2, nil
	}

	return nil, status.Error(codes.NotFound, "No user with that ID found.")
}
