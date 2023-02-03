package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

	defer cont.Terminate(ctx)

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

	conn := db.NewDbConnectionFromSettings(ctx, &dbset, true)
	conn.WaitForDB(logger)

	type Scenario struct {
		Name     string
		LastWeek []*models.Reward
		ThisWeek []*models.Reward
	}

	scens := []Scenario{
		{
			Name: "Level1Grow",
			LastWeek: []*models.Reward{
				{UserDeviceID: activeAutoPi, ConnectionStreak: 1, DisconnectionStreak: 0},
			},
			ThisWeek: []*models.Reward{
				{UserDeviceID: activeAutoPi, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
		},
		{
			Name: "Level1Disconnect",
			LastWeek: []*models.Reward{
				{UserDeviceID: inactive, ConnectionStreak: 1, DisconnectionStreak: 0},
			},
			ThisWeek: []*models.Reward{
				{UserDeviceID: inactive, ConnectionStreak: 1, DisconnectionStreak: 1, StreakPoints: 0, IntegrationPoints: 0},
			},
		},
		{
			Name: "JoinLevel2",
			LastWeek: []*models.Reward{
				{UserDeviceID: activeAutoPi, ConnectionStreak: 3, DisconnectionStreak: 0},
			},
			ThisWeek: []*models.Reward{
				{UserDeviceID: activeAutoPi, ConnectionStreak: 4, DisconnectionStreak: 0, StreakPoints: 1000, IntegrationPoints: 6000},
			},
		},
		{
			Name: "Level2Tesla",
			LastWeek: []*models.Reward{
				{UserDeviceID: activeTesla, ConnectionStreak: 5, DisconnectionStreak: 0},
			},
			ThisWeek: []*models.Reward{
				{UserDeviceID: activeTesla, ConnectionStreak: 6, DisconnectionStreak: 0, StreakPoints: 1000, IntegrationPoints: 4000},
			},
		},
		{
			Name: "LosingLevel",
			LastWeek: []*models.Reward{
				{UserDeviceID: inactive, ConnectionStreak: 22, DisconnectionStreak: 2},
			},
			ThisWeek: []*models.Reward{
				{UserDeviceID: inactive, ConnectionStreak: 4, DisconnectionStreak: 3, StreakPoints: 0, IntegrationPoints: 0},
			},
		},
		{
			Name:     "NewSmartcar",
			LastWeek: []*models.Reward{},
			ThisWeek: []*models.Reward{
				{UserDeviceID: activeSmartcar, ConnectionStreak: 1, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 1000},
			},
		},
	}

	for _, scen := range scens {
		t.Run(scen.Name, func(t *testing.T) {
			lastWk := models.IssuanceWeek{ID: 0, JobStatus: models.IssuanceWeeksJobStatusFinished}
			err = lastWk.Insert(ctx, conn.DBS().Writer, boil.Infer())
			if err != nil {
				t.Fatal(err)
			}

			for _, lst := range scen.LastWeek {
				err := lst.Insert(ctx, conn.DBS().Writer, boil.Infer())
				if err != nil {
					t.Fatal(err)
				}
			}

			mp := map[string]*models.Reward{}

			for _, ths := range scen.ThisWeek {
				mp[ths.UserDeviceID] = ths
			}

			task := RewardsTask{
				Logger:          &logger,
				DataService:     Views{},
				DefsClient:      &FakeDefClient{},
				DevicesClient:   &FakeDevClient{},
				DB:              conn,
				TransferService: &FakeTransfer{},
			}

			err = task.Calculate(1)
			if err != nil {
				t.Fatal(err)
			}

			rs, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(1), qm.OrderBy(models.RewardColumns.IssuanceWeekID+","+models.RewardColumns.UserDeviceID)).All(ctx, conn.DBS().Reader)
			if err != nil {
				t.Fatal(err)
			}

			mp2 := map[string]*models.Reward{}
			for _, ths := range rs {
				mp2[ths.UserDeviceID] = ths
			}

			for k, v1 := range mp {
				v2, ok := mp2[k]
				if !ok {
					t.Errorf("Missing row for device %s", k)
					continue
				}
				if v2.ConnectionStreak != v1.ConnectionStreak {
					t.Errorf("Device %s should have streak %d but had streak %d", k, v1.ConnectionStreak, v2.ConnectionStreak)
				}
				if v2.DisconnectionStreak != v1.DisconnectionStreak {
					t.Errorf("Device %s should have streak %d but had streak %d", k, v1.ConnectionStreak, v2.ConnectionStreak)
				}
				if v2.IntegrationPoints != v1.IntegrationPoints {
					t.Errorf("Device %s should have %d integration points but had %d", k, v1.IntegrationPoints, v2.IntegrationPoints)
				}
				if v2.StreakPoints != v1.StreakPoints {
					t.Errorf("Device %s should have %d streak points but had %d", k, v1.StreakPoints, v2.StreakPoints)
				}
			}

			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

			_, err = models.IssuanceWeeks().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}
		})

	}

}

type Views struct {
}

const autoPiIntegration = "2LFD6DXuGRdVucJO1a779kEUiYi"
const teslaIntegration = "2LFQOgsYd5MEmRNBnsYXKp0QHC3"
const smartcarIntegration = "2LFSA81Oo4agy0y4NvP7f6hTdgs"

const activeAutoPi = "2LFD2qeDxWMf49jSdEGQ2Znde3l"
const activeTesla = "2LFQTaaEzsUGyO2m1KtDIz4cgs0"
const activeSmartcar = "2LFSD4V6NcW88t3pdjPTNUJTPOu"
const inactive = "2LFOozehPU5ntHkuqHSQbn93seV"

func (v Views) DescribeActiveDevices(start, end time.Time) ([]*DeviceData, error) {
	return []*DeviceData{
		{ID: activeAutoPi, Integrations: []string{autoPiIntegration}},
		{ID: activeTesla, Integrations: []string{teslaIntegration}},
		{ID: activeSmartcar, Integrations: []string{smartcarIntegration}},
	}, nil
}

type FakeDefClient struct{}

func (d *FakeDefClient) GetIntegrations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_defs.GetIntegrationResponse, error) {
	return &pb_defs.GetIntegrationResponse{Integrations: []*pb_defs.Integration{
		{Id: autoPiIntegration, Vendor: "AutoPi"},
		{Id: teslaIntegration, Vendor: "Tesla"},
		{Id: smartcarIntegration, Vendor: "SmartCar"},
	}}, nil
}

type FakeDevClient struct{}

var rewards = map[string]*pb_devices.UserDevice{
	activeAutoPi: {
		Id:                       activeAutoPi,
		UserId:                   "USER1",
		TokenId:                  ref(uint64(1)),
		OptedInAt:                timestamppb.Now(),
		OwnerAddress:             common.FromHex("0x67B94473D81D0cd00849D563C94d0432Ac988B49"),
		AftermarketDeviceTokenId: ref(uint64(2)),
	},
	inactive: {},
	activeTesla: {
		Id:           activeTesla,
		UserId:       "USER1",
		TokenId:      ref(uint64(2)),
		OptedInAt:    timestamppb.Now(),
		OwnerAddress: common.FromHex("0x67B94473D81D0cd00849D563C94d0432Ac988B49"),
	},
	activeSmartcar: {
		Id:           activeSmartcar,
		UserId:       "USER1",
		TokenId:      ref(uint64(3)),
		OptedInAt:    timestamppb.Now(),
		OwnerAddress: common.FromHex("0x67B94473D81D0cd00849D563C94d0432Ac988B49"),
	},
}

func ref[A any](a A) *A {
	return &a
}

func (d *FakeDevClient) GetUserDevice(ctx context.Context, in *pb_devices.GetUserDeviceRequest, opts ...grpc.CallOption) (*pb_devices.UserDevice, error) {
	ud, ok := rewards[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "No device with that ID found.")
	}
	return ud, nil
}

type FakeTransfer struct{}

func (t *FakeTransfer) TransferUserTokens(ctx context.Context, week int) error {
	return nil
}
