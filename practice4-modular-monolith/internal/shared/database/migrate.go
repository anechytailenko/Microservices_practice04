package database

import (
	"database/sql"
	"embed"
	"log"

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
		log.Printf("[%s] Database is up to date. No new migrations to apply.", dirPath)
		return nil
	}

	if err != nil {
		return err
	}

	log.Printf("[%s] Migrations applied successfully", dirPath)
	return nil
}
