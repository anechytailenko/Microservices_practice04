package database

import (
	"context"
	"database/sql"
	"embed"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(db *sql.DB, fs embed.FS, dirPath string) error {
	sourceDriver, err := iofs.New(fs, dirPath)
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		logger.Printf(context.Background(), "[%s] Database is up to date. No new migrations to apply.", dirPath)
		return nil
	}

	if err != nil {
		return err
	}

	logger.Printf(context.Background(), "[%s] Migrations applied successfully", dirPath)
	return nil
}
