package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ericlagergren/decimal"

	"github.com/Shopify/sarama/mocks"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
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

type User struct {
	ID       string
	Address  common.Address
	Code     string
	CodeUsed string
}

var mkAddr = func(i int) common.Address {
	return common.BigToAddress(big.NewInt(int64(i)))
}

type Referral struct {
	Referee  common.Address
	Referrer common.Address
}

type FakeUserClient struct {
	users []User
}

var zeroAddr common.Address

func (d *FakeUserClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.User, error) {
	for _, user := range d.users {
		if user.ID == in.Id {
			addr := user.Address.Hex()
			out := &pb.User{
				Id:              user.ID,
				EthereumAddress: &addr,
			}

			if user.CodeUsed != "" {
				for _, ref := range d.users {
					if user.CodeUsed == ref.Code {
						if ref.Address != zeroAddr {
							out.ReferredBy = &pb.UserReferrer{
								EthereumAddress: ref.Address.Bytes(),
							}
						}
						break
					}
				}
			}

			return out, nil
		}
	}

	return nil, status.Error(codes.NotFound, "No user with that ID found.")
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

	type Device struct {
		ID      string
		TokenID int
		UserID  string
		VIN     string
	}

	type Reward struct {
		Week     int
		DeviceID string
		UserID   string
		Earning  bool
	}

	type Scenario struct {
		Name string
		// ReferralCount int
		// LastWeek      []*models.Reward
		// ThisWeek      []*models.Reward
		Users     []User
		Devices   []Device
		Rewards   []Reward
		Referrals []Referral
	}

	scens := []Scenario{
		{
			Name: "New address, new car, referred by non-deleted user",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, VIN: "00000000000000001"},
			},
			Users: []User{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2"},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{
				{Referee: mkAddr(1), Referrer: mkAddr(2)},
			},
		},
		{
			Name: "New address, new car, not referred",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, VIN: "00000000000000001"},
			},
			Users: []User{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{},
		},
		{
			Name: "Referrer has same address",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, VIN: "00000000000000001"},
			},
			Users: []User{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2"},
				{ID: "User2", Address: mkAddr(1), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{},
		},
		{
			Name: "Referrer has same address",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, VIN: "00000000000000001"},
			},
			Users: []User{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2"},
				{ID: "User2", Address: mkAddr(1), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{},
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

			for _, lst := range scen.Rewards {
				wk := models.IssuanceWeek{
					ID:        lst.Week,
					JobStatus: models.IssuanceWeeksJobStatusFinished,
				}
				wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())

				r := models.Reward{
					IssuanceWeekID: lst.Week,
					UserDeviceID:   lst.DeviceID,
					UserID:         lst.UserID,
				}
				if lst.Earning {
					r.Tokens = types.NewNullDecimal(decimal.New(100, 0))
				}

				err := r.Insert(ctx, conn.DBS().Writer, boil.Infer())
				if err != nil {
					t.Fatal(err)
				}
			}

			producer := mocks.NewSyncProducer(t, nil)
			transferService := NewTokenTransferService(&settings, producer, conn)

			referralBonusService := NewReferralBonusService(&settings, transferService, 1, &logger, &FakeUserClient{users: scen.Users})

			refs, err := referralBonusService.CollectReferrals(ctx, 5)
			if err != nil {
				t.Fatal(err)
			}

			var actual []Referral

			for i := 0; i < len(refs.Referees); i++ {
				actual = append(actual, Referral{Referee: refs.Referees[i], Referrer: refs.Referrers[i]})
			}

			assert.ElementsMatch(t, scen.Referrals, actual)
		})

	}

}

func TestReferralsBatchRequest(t *testing.T) {

	ctx := context.Background()

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	settings.TransferBatchSize = 1

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

	config := mocks.NewTestConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer := mocks.NewSyncProducer(t, config)

	transferService := NewTokenTransferService(&settings, producer, conn)

	referralBonusService := NewReferralBonusService(&settings, transferService, 1, &logger, &FakeUserClient{})

	refs := Referrals{
		Referees:  []common.Address{mkAddr(1)},
		Referrers: []common.Address{mkAddr(2)},
	}

	var out []shared.CloudEvent[transferData]

	checker := func(b2 []byte) error {
		var o shared.CloudEvent[transferData]
		fmt.Println("XFFCUTE", string(b2))
		err := json.Unmarshal(b2, &o)
		require.NoError(t, err)
		out = append(out, o)
		return nil
	}

	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	require.NoError(t, referralBonusService.transfer(ctx, refs))

	producer.Close()

	abi, err := contracts.ReferralMetaData.GetAbi()
	require.NoError(t, err)

	args := common.Bytes2Hex(abi.Methods["sendReferralBonuses"].ID)
	fmt.Println("SIGNATURE", args)
	for i, a := range abi.Methods["sendReferralBonuses"].Inputs {
		fmt.Println("A", i, a.Name, a.Type)
	}

	// type XL struct {
	// 	Referees  []common.Address
	// 	Referrers []common.Address
	// }

	var r []Referral

	fmt.Println("OUT", len(out))

	for i := range out {
		fmt.Println(i)
		xl := map[string]any{}
		fmt.Println("DDX", string(out[i].Data.Data))
		b := common.FromHex(out[i].Data.Data)
		fmt.Println("L1", len(b), common.Bytes2Hex(b))
		fmt.Println("L2", len(b[4:]), common.Bytes2Hex(b[4:]))
		require.NoError(t, err)
		fmt.Println("XPP", xl)
		o, _ := abi.Methods["sendReferralBonuses"].Inputs.Unpack(b[4:])
		referees := o[0].([]common.Address)
		referrers := o[1].([]common.Address)
		for i := range referees {
			r = append(r, Referral{Referee: referees[i], Referrer: referrers[i]})
		}
	}

	assert.ElementsMatch(t, []Referral{{Referee: mkAddr(1), Referrer: mkAddr(2)}}, r)
}
