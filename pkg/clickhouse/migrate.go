package clickhouse

import (
	"database/sql"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/prismelabs/analytics/pkg/log"
)

// migrate starts migrating a clickhouse instance to the latest version.
func migrate(logger log.Logger, db *sql.DB, source source.Driver) {
	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{
		MigrationsTable:       "migrations",
		MigrationsTableEngine: "MergeTree",
		MultiStatementEnabled: true,
	})
	if err != nil {
		logger.Panic().Msgf("failed to create golang-migrate driver for clickhouse migration: %v", err.Error())
	}

	m, err := gomigrate.NewWithInstance("migrations", source, "analytics", driver)
	if err != nil {
		logger.Panic().Msgf("failed to create go-migrate.Migrate instance: %v", err.Error())
	}
	m.Log = &logger

	err = m.Up()
	if err != nil && err != gomigrate.ErrNoChange {
		logger.Panic().Msgf("failed to execute clickhouse migrations: %v", err.Error())
	}

	logger.Info().Msg("clickhouse migration successfully done")
}
