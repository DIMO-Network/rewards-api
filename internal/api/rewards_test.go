package api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/shared/api/rewards"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
)

func TestGetBlacklist(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.New(os.Stdout)
	cont, dbs := utils.GetDbConnection(ctx, t, logger)
	defer testcontainers.CleanupContainer(t, cont)

	addr1 := common.HexToAddress("0x1")
	addr2 := common.HexToAddress("0x2")

	serv := NewRewardsService(dbs, &logger)

	creatTime := time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)

	timeNow = func() time.Time {
		return creatTime
	}

	_, err := serv.SetBlacklistStatus(ctx, &rewards.SetBlacklistStatusRequest{
		EthereumAddress: addr1.Bytes(),
		IsBlacklisted:   true,
		Note:            "xdd",
	})
	if err != nil {
		panic(err)
	}

	res, err := serv.GetBlacklistStatus(ctx, &rewards.GetBlacklistStatusRequest{
		EthereumAddress: addr1.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	if !res.IsBlacklisted {
		t.Errorf("%s should have been blacklisted", addr1)
	}

	if res.Note != "xdd" {
		t.Errorf("should return the note for %s", addr1)
	}

	if res.CreatedAt.AsTime() != creatTime {
		t.Errorf("wrong timestamp %s for creation", res.CreatedAt.AsTime())
	}

	res, err = serv.GetBlacklistStatus(ctx, &rewards.GetBlacklistStatusRequest{
		EthereumAddress: addr2.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	if res.IsBlacklisted {
		t.Errorf("should not have been ")
	}

	_, err = serv.SetBlacklistStatus(ctx, &rewards.SetBlacklistStatusRequest{
		EthereumAddress: addr1.Bytes(),
		IsBlacklisted:   false,
	})
	if err != nil {
		panic(err)
	}

	res, err = serv.GetBlacklistStatus(ctx, &rewards.GetBlacklistStatusRequest{
		EthereumAddress: addr1.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	if res.IsBlacklisted {
		t.Errorf("should have been removed")
	}
}
