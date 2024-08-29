package database

import (
	"context"
	"database/sql"

	"github.com/DIMO-Network/shared/db"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

func MigrateDatabase(_ zerolog.Logger, settings *db.Settings, command string, dir string) error {
	db, err := sql.Open("postgres", settings.BuildConnectionString(true))
	if err != nil {
		return err
	}
	defer db.Close() //nolint

	if err := db.Ping(); err != nil {
		return err
	}

	if command == "" {
		command = "up"
	}

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS rewards_api;"); err != nil {
		return err
	}

	goose.SetTableName("rewards_api.migrations")

	return goose.RunContext(context.TODO(), command, db, dir)
}
