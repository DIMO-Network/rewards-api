package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetDbConnection(ctx context.Context, t *testing.T, logger zerolog.Logger) (testcontainers.Container, db.Store) {
	port := 5432
	nport := fmt.Sprintf("%d/tcp", port)

	cont, err := postgres.Run(
		ctx,
		"docker.io/postgres:16.6-alpine",
		postgres.WithDatabase("rewards_api"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

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
	pdb := db.NewDbConnectionForTest(ctx, &dbset, false)
	for !pdb.IsReady() {
		time.Sleep(200 * time.Millisecond)
	}

	if err := database.MigrateDatabase(logger, &dbset, "", "../../migrations"); err != nil {
		t.Fatal(err)
	}

	conn := db.NewDbConnectionForTest(ctx, &dbset, true)
	conn.WaitForDB(logger)

	return cont, conn
}
