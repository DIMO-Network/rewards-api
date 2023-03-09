package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const newUserDeviceID = "2LFD2qeDxWMf49jSdEGQ2Znde3l"
const existingUserDeviceID = "2LFQTaaEzsUGyO2m1KtDIz4cgs0"

const newUserReferred = "NewUserReferred"
const newUserNotReferred = "newUserNotReferred"
const userDeletedTheirAccount = "userDeletedTheirAccount"
const existingUser = "ExistingUser"

var existingUserAddr = "0x41A82213D81D0cd00849D563C94d0432Ac988COL"
var newUserReferredAddr = "0x67B94473D81D0cd00849D563C94d0432Ac988B49"
var newUserNotReferredAddr = "0x67B94473D81D0cd00849D563C94d0432Ac988B48"
var userDeletedAccountAddr = "0x67B94473D81D0cd00849D563C94d0432Ac988B47"

var referrer = "0x67B94473D81D0cd00849D563C94d0432Ac988B50"

var fakeUserClientResponse = map[string]*pb_users.User{
	newUserReferred: {
		Id:              newUserReferred,
		EthereumAddress: &newUserReferredAddr,
		ReferredBy:      &pb_users.UserReferrer{EthereumAddress: common.FromHex(referrer)},
	},
	newUserNotReferred: {
		Id:              newUserReferred,
		EthereumAddress: &newUserNotReferredAddr,
		ReferredBy:      nil,
	},
	userDeletedTheirAccount: {
		Id:              userDeletedTheirAccount,
		EthereumAddress: &userDeletedAccountAddr,
	},
}

type FakeUserClient struct{}

func (d *FakeUserClient) GetUser(ctx context.Context, in *pb_users.GetUserRequest, opts ...grpc.CallOption) (*pb_users.User, error) {
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
		t.Fatalf("failed to start geneic container: %v", err)
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
		Name         string
		Referral     bool
		ReferredBy   *pb_users.UserReferrer
		NewUserCount int
		LastWeek     []*models.Reward
		ThisWeek     []*models.Reward
	}

	scens := []Scenario{
		{
			Name:         newUserReferred,
			Referral:     true,
			ReferredBy:   &pb_users.UserReferrer{EthereumAddress: common.FromHex(referrer)},
			NewUserCount: 1,
			LastWeek: []*models.Reward{
				{UserID: existingUser, UserEthereumAddress: null.StringFrom(existingUserAddr), IssuanceWeekID: 0, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
			ThisWeek: []*models.Reward{
				{UserID: existingUser, UserEthereumAddress: null.StringFrom(existingUserAddr), IssuanceWeekID: 1, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
				{UserID: newUserReferred, UserEthereumAddress: null.StringFrom(newUserReferredAddr), IssuanceWeekID: 1, UserDeviceID: newUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
		},
		{
			Name:         newUserNotReferred,
			Referral:     true,
			ReferredBy:   nil,
			NewUserCount: 0,
			LastWeek: []*models.Reward{
				{UserID: existingUser, UserEthereumAddress: null.StringFrom(existingUserAddr), IssuanceWeekID: 0, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
			ThisWeek: []*models.Reward{
				{UserID: existingUser, UserEthereumAddress: null.StringFrom(existingUserAddr), IssuanceWeekID: 1, UserDeviceID: existingUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
				{UserID: newUserNotReferred, UserEthereumAddress: null.StringFrom(newUserNotReferredAddr), IssuanceWeekID: 1, UserDeviceID: newUserDeviceID, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
		},
		{
			Name:         userDeletedTheirAccount,
			Referral:     true,
			NewUserCount: 0,
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

			producer := mocks.NewSyncProducer(t, nil)
			transferService := NewTokenTransferService(&settings, producer, conn)

			referralBonusService := NewReferralBonusService(&settings, transferService, 1, nil, &FakeUserClient{})

			weeklyRefs, err := referralBonusService.CollectReferrals(ctx, 1)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("name: ", scen.Name)
			assert.Equal(t, len(weeklyRefs.Referrees), scen.NewUserCount)
			assert.Equal(t, len(weeklyRefs.Referrers), scen.NewUserCount)
			for _, addr := range weeklyRefs.Referrers {
				userClientResp := fakeUserClientResponse[scen.Name]
				userReferredBy := *userClientResp.GetReferredBy()
				assert.Equal(t, common.BytesToAddress(userReferredBy.GetEthereumAddress()), common.BytesToAddress(scen.ReferredBy.EthereumAddress))
				assert.Equal(t, addr, common.BytesToAddress(scen.ReferredBy.EthereumAddress))
			}

		})

	}

}

func TestReferralsBatchRequest(t *testing.T) {

	config := mocks.NewTestConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer := mocks.NewSyncProducer(t, config)

	refs := Referrals{
		Referrees: []common.Address{common.HexToAddress("0x67B94473D81D0cd00849D563C94d0432Ac988B49")},
		Referrers: []common.Address{common.HexToAddress("0x67B94473D81D0cd00849D563C94d0432Ac988B48")},
	}

	abi, err := contracts.ReferralsMetaData.GetAbi()
	assert.Nil(t, err)

	data, err := abi.Pack("sendReferralBonuses", refs.Referrees, refs.Referrers)
	assert.Nil(t, err)

	event := shared.CloudEvent[string]{
		ID:          "",
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     "contract addr",
		Time:        time.Now(),
		Type:        "zone.dimo.referrals.transaction.request",
		Data:        hexutil.Encode(data),
	}

	eventBytes, err := json.Marshal(event)
	assert.Nil(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, eventBytes, b2)
		return nil
	}

	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	if _, _, err = producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: "test",
			Key:   sarama.StringEncoder(""),
			Value: sarama.ByteEncoder(eventBytes),
		},
	); err != nil {
		assert.Nil(t, err)
	}

}
