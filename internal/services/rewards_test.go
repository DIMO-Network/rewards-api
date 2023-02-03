package services

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetWeekNumForCron(t *testing.T) {
	ti, _ := time.Parse(time.RFC3339, "2022-02-07T05:00:02Z")
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}

	ti, _ = time.Parse(time.RFC3339, "2022-02-07T04:58:44Z")
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}
}

func TestCalc(t *testing.T) {
	ctx := context.Background()

	port := "5432/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12.11-alpine",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
		},
		WaitingFor: wait.ForListeningPort(nat.Port(port)),
	}
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// connStr :=
}
