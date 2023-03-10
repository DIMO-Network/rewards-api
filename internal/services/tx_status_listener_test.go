package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/docker/go-connections/nat"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func TestBaselineStatus(t *testing.T) {
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

	dbs := db.NewDbConnectionForTest(ctx, &dbset, true)
	dbs.WaitForDB(logger)

	proc, err := NewStatusProcessor(dbs, &logger)
	if err != nil {
		t.Fatal(err)
	}

	iw := models.IssuanceWeek{ID: 5, JobStatus: models.IssuanceWeeksJobStatusFinished}
	err = iw.Insert(ctx, dbs.DBS().Writer, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	reqID := ksuid.New().String()

	mtr := models.MetaTransactionRequest{ID: reqID, Status: models.MetaTransactionRequestStatusMined}
	err = mtr.Insert(ctx, dbs.DBS().Writer, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	r := models.Reward{IssuanceWeekID: iw.ID, UserDeviceID: ksuid.New().String(), UserID: "XDD", ConnectionStreak: 3, DisconnectionStreak: 0, TransferMetaTransactionRequestID: null.StringFrom(reqID), UserDeviceTokenID: types.NewNullDecimal(decimal.New(55, 0))}

	err = r.Insert(ctx, dbs.DBS().Writer, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	x := fmt.Sprintf(`{
	"data": {
		"requestId": %q,
		"type": "Confirmed",
		"transaction": {
			"successful": true,
			"logs": [
				{
					"topics": [
						"0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea",
						"0x000000000000000000000000e4da3218b897e3f72ada9f5cabc2c9d61983bd92",
						"0x0000000000000000000000000000000000000000000000000000000000000037"
					],
					"data": "0x0000000000000000000000000000000000000000000000000000000008f0d180"
				}
			]
		}
	}
}`, reqID)

	err = proc.processMessage(&sarama.ConsumerMessage{Value: []byte(x)})
	if err != nil {
		t.Fatal(err)
	}

	err = mtr.Reload(ctx, dbs.DBS().Reader)
	if err != nil {
		t.Fatal(err)
	}

	if mtr.Status != models.MetaTransactionRequestStatusConfirmed {
		t.Error("expected transaction to fall into Confirmed")
	}

	if mtr.Successful != null.BoolFrom(true) {
		t.Error("expected transaction to be marked as successful")
	}

	err = r.Reload(ctx, dbs.DBS().Reader)
	if err != nil {
		t.Fatal(err)
	}

	if r.TransferSuccessful != null.BoolFrom(true) {
		t.Error("not set as successful")
	}
}
