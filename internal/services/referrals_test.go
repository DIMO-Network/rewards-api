package services

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ericlagergren/decimal"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/IBM/sarama/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
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

func TestReferrals(t *testing.T) {
	ctx := context.Background()

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		t.Fatal(err)
	}

	refContractAddr := common.HexToAddress("0xfF358a3dB687d9E80435a642bB3Ba8E64D4359A6")
	logger := zerolog.New(os.Stdout)
	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	type Device struct {
		ID               string
		TokenID          int
		UserID           string
		Vin              string
		FirstEarningWeek int
	}

	type Reward struct {
		Week             int
		DeviceID         string
		UserID           string
		Earning          bool
		ConnectionStreak int
	}

	type Scenario struct {
		Name      string
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
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 1},
				{Week: 6, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 2},
				{Week: 7, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 3},
				{Week: 8, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 4},
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
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 1},
				{Week: 6, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 2},
				{Week: 7, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 3},
				{Week: 8, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 4},
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
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 1},
				{Week: 6, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 2},
				{Week: 7, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 3},
				{Week: 8, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 4},
			},
			Referrals: []Referral{},
		},
		{
			Name: "Referring user was deleted",
			Devices: []Device{
				{ID: "Dev1", UserID: "User1", TokenID: 1, Vin: "00000000000000001", FirstEarningWeek: 5},
			},
			Users: []refUser{
				{ID: "User1", Address: mkAddr(1), Code: "1", CodeUsed: "2", InvalidReferrer: true},
				{ID: "User2", Address: mkAddr(2), Code: "2", CodeUsed: ""},
			},
			Rewards: []Reward{
				{Week: 5, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 1},
				{Week: 6, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 2},
				{Week: 7, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 3},
				{Week: 8, DeviceID: "Dev1", UserID: "User1", Earning: true, ConnectionStreak: 4},
			},
			Referrals: []Referral{
				{Referee: mkAddr(1), Referrer: refContractAddr},
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
					IssuanceWeekID:   lst.Week,
					UserDeviceID:     lst.DeviceID,
					UserID:           lst.UserID,
					ConnectionStreak: lst.ConnectionStreak,
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

			ctrl := gomock.NewController(t)
			usersClient := NewMockUsersClient(ctrl)
			accountsClient := NewMockAccountsClient(ctrl)
			accountsClient.EXPECT().TempReferral(gomock.Any(), gomock.Any()).Return(
				nil, status.Error(codes.NotFound, "Not found xdd"),
			)

			usersClient.EXPECT().GetUsersByEthereumAddress(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, in *pb.GetUsersByEthereumAddressRequest, opts ...grpc.CallOption) (*pb.GetUsersByEthereumAddressResponse, error) {
					for _, u := range scen.Users {
						if u.Address == common.BytesToAddress(in.EthereumAddress) {
							uOut := &pb.User{Id: u.ID, EthereumAddressBytes: u.Address.Bytes()}

							if u.InvalidReferrer {
								uOut.ReferredBy = &pb.UserReferrer{
									ReferrerValid: false,
								}
							} else if u.CodeUsed != "" {
								for _, u2 := range scen.Users {
									if u2.Code == u.CodeUsed {
										uOut.ReferredBy = &pb.UserReferrer{
											Id:              u2.ID,
											EthereumAddress: u2.Address.Bytes(),
											ReferrerValid:   true,
										}
										break
									}
								}
							}

							return &pb.GetUsersByEthereumAddressResponse{Users: []*pb.User{uOut}}, nil
						}
					}
					return &pb.GetUsersByEthereumAddressResponse{Users: []*pb.User{}}, nil
				}).AnyTimes()

			referralBonusService := NewReferralBonusService(&settings, transferService, 1, &logger, usersClient, accountsClient)
			referralBonusService.ContractAddress = refContractAddr

			refs, err := referralBonusService.CollectReferrals(ctx, 8)
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

	logger := zerolog.New(os.Stdout)
	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	config := mocks.NewTestConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer := mocks.NewSyncProducer(t, config)

	transferService := NewTokenTransferService(&settings, producer, conn)

	ctrl := gomock.NewController(t)
	usersClient := NewMockUsersClient(ctrl)
	accountsClient := NewMockAccountsClient(ctrl)

	referralBonusService := NewReferralBonusService(&settings, transferService, 1, &logger, usersClient, accountsClient)

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
