package services

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID       string
	Address  common.Address
	Code     string
	CodeUsed string
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
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}

	ti, _ = time.Parse(time.RFC3339, "2022-02-07T04:58:44Z")
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}
}

func TestStreak(t *testing.T) {
	ctx := context.Background()

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	port := 5432
	nport := fmt.Sprintf("%d/tcp", port)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:12.11-alpine",
		ExposedPorts: []string{nport},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_DB":       "rewards_api",
			"POSTGRES_PASSWORD": "postgres",
		},
		WaitingFor: wait.ForListeningPort(nat.Port(nport)),
	}
	cont, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	defer cont.Terminate(ctx) //nolint

	logger := zerolog.Nop()

	host, err := cont.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mport, err := cont.MappedPort(ctx, nat.Port(nport))
	if err != nil {
		t.Fatal(err)
	}

	dbset := db.Settings{
		User:               "postgres",
		Password:           "postgres",
		Port:               mport.Port(),
		Host:               host,
		Name:               "rewards_api",
		MaxOpenConnections: 10,
		MaxIdleConnections: 10,
	}

	if err := database.MigrateDatabase(logger, &dbset, "", "../../migrations"); err != nil {
		t.Fatal(err)
	}

	conn := db.NewDbConnectionForTest(ctx, &dbset, true)
	conn.WaitForDB(logger)

	type OldReward struct {
		Week       int
		DeviceID   string
		UserID     string
		ConnStreak int
		DiscStreak int
	}

	type NewReward struct {
		DeviceID          string
		TokenID           int
		Address           common.Address
		ConnStreak        int
		DiscStreak        int
		StreakPoints      int
		IntegrationPoints int
	}

	type VIN struct {
		VIN        string
		FirstWeek  int
		FirstToken int
	}

	type Scenario struct {
		Name     string
		Previous []OldReward
		Devices  []Device
		Users    []User
		New      []NewReward
		PrevVIN  []VIN
		NewVIN   []VIN
	}

	scens := []Scenario{
		{
			Name: "Level1SmartcarGrow",
			Users: []User{
				{ID: "User1", Address: mkAddr(1)},
			},
			Devices: []Device{
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{smartcarIntegration}, Opted: true},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 2, DiscStreak: 0, StreakPoints: 0, IntegrationPoints: 1000},
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
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}, Opted: true},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 1, DiscStreak: 0},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 1, StreakPoints: 0, IntegrationPoints: 0},
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
				{ID: mkID(1), TokenID: 1, AMID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{autoPiIntegration}, Opted: true},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 3, DiscStreak: 0},
			},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 4, DiscStreak: 0, StreakPoints: 1000, IntegrationPoints: 6000},
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
				{ID: mkID(1), TokenID: 1, AMID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}, Opted: true},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), UserID: "User1", ConnStreak: 22, DiscStreak: 2},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 4, DiscStreak: 3, StreakPoints: 0, IntegrationPoints: 0},
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
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{teslaIntegration}, Opted: true},
			},
			Previous: []OldReward{},
			New: []NewReward{
				{DeviceID: mkID(1), TokenID: 1, Address: mkAddr(1), ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, IntegrationPoints: 4000},
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
				{ID: mkID(1), TokenID: 1, UserID: "User1", VIN: mkVIN(1), IntsWithData: []string{}, Opted: true},
				{ID: mkID(2), TokenID: 2, UserID: "User2", VIN: mkVIN(1), IntsWithData: []string{teslaIntegration}, Opted: true},
			},
			Previous: []OldReward{
				{Week: 4, DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 0},
			},
			New: []NewReward{
				{DeviceID: mkID(1), ConnStreak: 1, DiscStreak: 1, StreakPoints: 0, IntegrationPoints: 0},
				{DeviceID: mkID(2), TokenID: 2, Address: mkAddr(2), ConnStreak: 1, DiscStreak: 0, StreakPoints: 0, IntegrationPoints: 4000},
			},
			PrevVIN: []VIN{{VIN: mkVIN(1), FirstToken: 1}},
			NewVIN:  []VIN{},
		},
	}

	for _, scen := range scens {
		t.Run(scen.Name, func(t *testing.T) {
			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

			models.Vins().DeleteAll(ctx, conn.DBS().Writer)

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

				wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())

				rw := models.Reward{
					IssuanceWeekID:      lst.Week,
					UserDeviceID:        lst.DeviceID,
					ConnectionStreak:    lst.ConnStreak,
					DisconnectionStreak: lst.DiscStreak,
				}
				err := rw.Insert(ctx, conn.DBS().Writer, boil.Infer())
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

				vw.Insert(ctx, conn.DBS().Writer.DB, boil.Infer())
			}

			transferService := NewTokenTransferService(&settings, nil, conn)

			rwBonusService := NewBaselineRewardService(&settings, transferService, Views{devices: scen.Devices}, &FakeDevClient{devices: scen.Devices, users: scen.Users}, &FakeDefClient{}, 5, &logger)

			rwBonusService.Calculate(5)

			rw, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(5), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			var actual []NewReward

			for _, c := range rw {
				nr := NewReward{
					DeviceID:          c.UserDeviceID,
					ConnStreak:        c.ConnectionStreak,
					DiscStreak:        c.DisconnectionStreak,
					StreakPoints:      c.StreakPoints,
					IntegrationPoints: c.IntegrationPoints,
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

			vs, err := models.Vins(models.VinWhere.FirstEarningWeek.EQ(5)).All(ctx, conn.DBS().Reader.DB)
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

type Device struct {
	ID           string
	TokenID      int
	UserID       string
	VIN          string
	Opted        bool
	IntsWithData []string
	AMID         int
}

type Views struct {
	devices []Device
}

const autoPiIntegration = "2LFD6DXuGRdVucJO1a779kEUiYi"
const teslaIntegration = "2LFQOgsYd5MEmRNBnsYXKp0QHC3"
const smartcarIntegration = "2LFSA81Oo4agy0y4NvP7f6hTdgs"

func (v Views) DescribeActiveDevices(start, end time.Time) ([]*DeviceData, error) {
	var out []*DeviceData
	for _, d := range v.devices {
		if len(d.IntsWithData) == 0 {
			continue
		}
		out = append(out, &DeviceData{
			ID: d.ID, Integrations: d.IntsWithData,
		})
	}
	return out, nil
}

type FakeDefClient struct {
}

func (d *FakeDefClient) GetIntegrations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_defs.GetIntegrationResponse, error) {
	return &pb_defs.GetIntegrationResponse{Integrations: []*pb_defs.Integration{
		{Id: autoPiIntegration, Vendor: "AutoPi"},
		{Id: teslaIntegration, Vendor: "Tesla"},
		{Id: smartcarIntegration, Vendor: "SmartCar"},
	}}, nil
}

type FakeDevClient struct {
	users   []User
	devices []Device
}

var zeroAddr common.Address

func (d *FakeDevClient) GetUserDevice(ctx context.Context, in *pb_devices.GetUserDeviceRequest, opts ...grpc.CallOption) (*pb_devices.UserDevice, error) {
	for _, ud := range d.devices {
		if ud.ID != in.Id {
			continue
		}

		var tk *uint64
		if ud.TokenID != 0 {
			t := uint64(ud.TokenID)
			tk = &t
		}

		var t *timestamppb.Timestamp
		if ud.Opted {
			t = timestamppb.New(time.Now())
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
			OptedInAt:    t,
			Vin:          vin,
			OwnerAddress: owner,
		}

		if ud.AMID != 0 {
			u1 := uint64(ud.AMID)
			ud2.AftermarketDeviceTokenId = &u1
		}

		return ud2, nil
	}

	return nil, status.Error(codes.NotFound, "No user with that ID found.")
}

type FakeTransfer struct{}

func (t *FakeTransfer) BaselineIssuance(ctx context.Context, week int) error {
	return nil
}
