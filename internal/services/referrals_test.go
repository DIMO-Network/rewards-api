package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/docker/go-connections/nat"
	"gotest.tools/assert"

	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const newUser = "2LFD2qeDxWMf49jSdEGQ2Znde3l"
const existingUser = "2LFQTaaEzsUGyO2m1KtDIz4cgs0"

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

	conn := db.NewDbConnectionFromSettings(ctx, &dbset, true)
	conn.WaitForDB(logger)

	type Scenario struct {
		Name         string
		Referral     bool
		NewUserCount int
		LastWeek     []*models.Reward
		ThisWeek     []*models.Reward
	}

	scens := []Scenario{
		{
			Name:         "newUser",
			Referral:     true,
			NewUserCount: 1,
			LastWeek: []*models.Reward{
				{UserID: "User1", IssuanceWeekID: 0, UserDeviceID: existingUser, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
			},
			ThisWeek: []*models.Reward{
				{UserID: "User2", IssuanceWeekID: 1, UserDeviceID: newUser, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
				{UserID: "User1", IssuanceWeekID: 1, UserDeviceID: existingUser, ConnectionStreak: 2, DisconnectionStreak: 0, StreakPoints: 0, IntegrationPoints: 6000},
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
				Logger: &logger,
				DB:     conn,
			}

			newUsrAddrs, err := task.CollectReferrals(1)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(newUsrAddrs), scen.NewUserCount)

		})

	}

}
