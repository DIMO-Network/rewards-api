package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ericlagergren/decimal"

	"github.com/Shopify/sarama/mocks"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type refUser struct {
	ID              string
	Address         common.Address
	Code            string
	CodeUsed        string
	InvalidReferrer bool
}

type Referral struct {
	Referee  common.Address
	Referrer common.Address
}

type FakeUserClient struct {
	users []refUser
}

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
								ReferrerValid:   true,
							}
						}

						if user.InvalidReferrer {
							out.ReferredBy.ReferrerValid = false
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

	refContractAddr := common.HexToAddress("0xfF358a3dB687d9E80435a642bB3Ba8E64D4359A6")

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
		ID               string
		TokenID          int
		UserID           string
		Vin              string
		FirstEarningWeek int
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
		Users     []refUser
		Devices   []Device
		Rewards   []Reward
		Referrals []Referral
	}

	scens := []Scenario{
		{
			Name: "New address, new car, referred by non-deleted user",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
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
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
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
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2"},
				{ID: "User2", Address: mkAddr(1), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{},
		},
		{
			Name: "Referring user has invalid wallet address",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2", InvalidReferrer: true},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
			},
			Referrals: []Referral{
				{Referee: mkAddr(1), Referrer: refContractAddr},
			},
		},
		{
			Name: "New address, new token, old Vin",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 0},
				{ID: "Dev3", UserID: "User3", TokenID: 3, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: ""},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
				{ID: "User3", Address: mkAddr(3), Code: "3", CodeUsed: "2"},
			},
			Rewards: []Reward{
				{Week: 3, DeviceID: "Dev1", UserID: "User1", Earning: true},
				{Week: 5, DeviceID: "Dev3", UserID: "User3", Earning: true},
			},
			Referrals: []Referral{},
		},
		{
			Name: "New Vin and user, same address",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
				{ID: "Dev2", UserID: "User2", TokenID: 3, Vin: "00000000000000002", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: ""},
				{ID: "User2", Address: mkAddr(1), Code: "2", CodeUsed: "3"},
				{ID: "User3", Address: mkAddr(3), Code: "3", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 3, DeviceID: "Dev1", UserID: "User1", Earning: true},
				{Week: 5, DeviceID: "Dev2", UserID: "User2", Earning: true},
			},
			Referrals: []Referral{},
		},
		{
			Name: "New user joins and connects a car that has previously been connected to DIMO",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 3},
				{ID: "Dev2", UserID: "User2", TokenID: 2, Vin: "00000000000000002", FirstEarningWeek: 3},
				{ID: "Dev2", UserID: "User3", TokenID: 3, Vin: "00000000000000002", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: ""},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
				{ID: "User3", Address: mkAddr(3), Code: "3", CodeUsed: "1"},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
				{Week: 5, DeviceID: "Dev2", UserID: "User3", Earning: true},
			},
			Referrals: []Referral{
				{Referee: mkAddr(3), Referrer: mkAddr(1)},
			},
		},
		{
			Name: "New user, two vehicles, only one genuinely new",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 3},
				{ID: "Dev2", UserID: "User2", TokenID: 2, Vin: "00000000000000002", FirstEarningWeek: 3},
				{ID: "Dev3", UserID: "User2", TokenID: 3, Vin: "00000000000000003", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: ""},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: "1"},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
				{Week: 5, DeviceID: "Dev2", UserID: "User2", Earning: true},
				{Week: 5, DeviceID: "Dev3", UserID: "User2", Earning: true},
			},
			Referrals: []Referral{
				{Referee: mkAddr(2), Referrer: mkAddr(1)},
			},
		},
		{
			Name: "Users who did not earn in a given week can still refer",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
				{ID: "Dev2", UserID: "User2", TokenID: 1, Vin: "00000000000000002", FirstEarningWeek: 1},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2"},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true},
				{Week: 5, DeviceID: "Dev2", UserID: "User2", Earning: false},
			},
			Referrals: []Referral{
				{Referee: mkAddr(1), Referrer: mkAddr(2)},
			},
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
				err := wk.Upsert(ctx, conn.DBS().Writer, false, []string{models.IssuanceWeekColumns.ID}, boil.Infer(), boil.Infer())
				assert.NoError(t, err)

				r := models.Reward{
					IssuanceWeekID: lst.Week,
					UserDeviceID:   lst.DeviceID,
					UserID:         lst.UserID,
				}
				if lst.Earning {
					r.Tokens = types.NewNullDecimal(decimal.New(100, 0))
				}
				for _, u := range scen.Users {
					if u.ID == lst.UserID {
						r.UserEthereumAddress = null.StringFrom(u.Address.Hex())
					}
				}
				for _, d := range scen.Devices {
					if d.ID == lst.DeviceID {
						r.UserDeviceTokenID = types.NewNullDecimal(decimal.New(int64(d.TokenID), 0))
					}
				}

				err = r.Insert(ctx, conn.DBS().Writer, boil.Infer())
				if err != nil {
					t.Fatal(err)
				}
			}

			producer := mocks.NewSyncProducer(t, nil)
			transferService := NewTokenTransferService(&settings, producer, conn)

			for _, ud := range scen.Devices {
				if ud.Vin != "" && len(ud.Vin) == 17 {
					vinRec := models.Vin{
						Vin:                 ud.Vin,
						FirstEarningWeek:    ud.FirstEarningWeek,
						FirstEarningTokenID: types.NewDecimal(new(decimal.Big).SetUint64(uint64(ud.TokenID))),
					}
					if err := vinRec.Upsert(ctx, transferService.db.DBS().Writer, false, []string{models.VinColumns.Vin}, boil.Infer(), boil.Infer()); err != nil {
						require.NoError(t, err)
					}
				}
			}

			referralBonusService := NewReferralBonusService(&settings, transferService, 1, &logger, &FakeUserClient{users: scen.Users})
			referralBonusService.ContractAddress = refContractAddr

			refs, err := referralBonusService.CollectReferrals(ctx, 5)
			require.NoError(t, err)

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
		Referees:        []common.Address{mkAddr(1)},
		Referrers:       []common.Address{mkAddr(2)},
		RefereeUserIDs:  []string{"xdd"},
		ReferrerUserIDs: []string{""},
	}

	var out []shared.CloudEvent[transferData]

	checker := func(b2 []byte) error {
		var o shared.CloudEvent[transferData]
		err := json.Unmarshal(b2, &o)
		require.NoError(t, err)
		out = append(out, o)
		return nil
	}

	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	wk := models.IssuanceWeek{
		ID:        referralBonusService.Week,
		JobStatus: models.IssuanceWeeksJobStatusFinished,
	}
	err = wk.Insert(ctx, conn.DBS().Writer, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	require.NoError(t, referralBonusService.transfer(ctx, refs))

	producer.Close()

	abi, err := contracts.ReferralMetaData.GetAbi()
	require.NoError(t, err)
	var r []Referral

	for i := range out {

		b := out[i].Data.Data
		require.NoError(t, err)
		o, _ := abi.Methods["sendReferralBonuses"].Inputs.Unpack(b[4:])
		referees := o[0].([]common.Address)
		referrers := o[1].([]common.Address)
		for i := range referees {
			r = append(r, Referral{Referee: referees[i], Referrer: referrers[i]})
		}
	}

	assert.ElementsMatch(t, []Referral{{Referee: mkAddr(1), Referrer: mkAddr(2)}}, r)
}
