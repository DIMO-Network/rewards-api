package api

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/shared/api/rewards"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetBlacklist(t *testing.T) {
	ctx := context.TODO()
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

	logger := zerolog.New(os.Stdout)

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

	dbs := db.NewDbConnectionForTest(ctx, &dbset, true)
	dbs.WaitForDB(logger)

	addr1 := common.HexToAddress("0x1")
	addr2 := common.HexToAddress("0x2")

	serv := NewRewardsService(dbs, &logger)

	_, err = serv.SetBlacklistStatus(ctx, &rewards.SetBlacklistStatusRequest{
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
		t.Errorf("should have been ")
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
