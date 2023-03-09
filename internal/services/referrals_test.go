package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/assert"

	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const newUserDeviceID = "2LFD2qeDxWMf49jSdEGQ2Znde3l"
const existingUserDeviceID = "2LFQTaaEzsUGyO2m1KtDIz4cgs0"

const newUserReferred = "NewUserReferred"
const newUserNotReferred = "newUserNotReferred"
const userDeletedTheirAccount = "userDeletedTheirAccount"
const existingUser = "ExistingUser"

var addr = "0x67B94473D81D0cd00849D563C94d0432Ac988B49"
var fakeUserClientResponse = map[string]*pb.User{
	newUserReferred: {
		Id:              newUserReferred,
		EthereumAddress: &addr,
		ReferredBy:      &pb.UserReferrer{EthereumAddress: common.FromHex("0x67B94473D81D0cd00849D563C94d0432Ac988B50")},
	},
	newUserNotReferred: {
		Id:              newUserReferred,
		EthereumAddress: &addr,
	},
	userDeletedTheirAccount: {
		Id:              userDeletedTheirAccount,
		EthereumAddress: &addr,
	},
}

type FakeUserClient struct{}

func (d *FakeUserClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.User, error) {
	ud, ok := fakeUserClientResponse[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "No user with that ID found.")
	}
	if ud.Id == userDeletedTheirAccount {
		return nil, nil
	}
	return ud, nil
}

func TestReferrals(t *testing.T) {
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

	conn := db.NewDbConnectionForTest(ctx, &dbset, true)
	conn.WaitForDB(logger)

	type Scenario struct {
		Name          string
		ReferralCount int
		LastWeek      []*models.Reward
		ThisWeek      []*models.Reward
	}

	scens := []Scenario{
		{
			Name:          newUserReferred,
			ReferralCount: 1,
			LastWeek: []*models.Reward{
				{UserID: existingUser, IssuanceWeekID: 0, UserDeviceID: existingUserDeviceID, ConnectionStreak: 1, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
			ThisWeek: []*models.Reward{
				{UserID: existingUser, IssuanceWeekID: 1, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
				{UserID: newUserReferred, IssuanceWeekID: 1, UserDeviceID: newUserDeviceID, ConnectionStreak: 1, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
		},
		{
			Name:          newUserNotReferred,
			ReferralCount: 0,
			LastWeek: []*models.Reward{
				{UserID: existingUser, IssuanceWeekID: 0, UserDeviceID: existingUserDeviceID, ConnectionStreak: 1, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
			ThisWeek: []*models.Reward{
				{UserID: existingUser, IssuanceWeekID: 1, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
				{UserID: newUserNotReferred, IssuanceWeekID: 1, UserDeviceID: newUserDeviceID, ConnectionStreak: 1, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
		},
		{
			Name:          userDeletedTheirAccount,
			ReferralCount: 0,
		},
	}

	for _, scen := range scens {
		t.Run(scen.Name, func(t *testing.T) {
			_, err = models.Rewards().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

			_, err = models.IssuanceWeeks().DeleteAll(ctx, conn.DBS().Writer)
			if err != nil {
				t.Fatal(err)
			}

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

			thisWk := models.IssuanceWeek{ID: 1, JobStatus: models.IssuanceWeeksJobStatusFinished}
			err = thisWk.Insert(ctx, conn.DBS().Writer, boil.Infer())
			if err != nil {
				t.Fatal(err)
			}

			for _, lst := range scen.ThisWeek {
				err := lst.Insert(ctx, conn.DBS().Writer, boil.Infer())
				if err != nil {
					t.Fatal(err)
				}
			}

			task := ReferralsTask{
				Logger:      &logger,
				DB:          conn,
				UsersClient: &FakeUserClient{},
			}

			weeklyRefs, err := task.CollectReferrals(ctx, 1)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(weeklyRefs.Referees), scen.ReferralCount)
			assert.Equal(t, len(weeklyRefs.Referrer), scen.ReferralCount)
		})
	}
}
