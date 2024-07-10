package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/Vikot10/poact/internal/config"
)

func main() {
	cfg := config.MustLoad()

	logger := mustCreateLogger(cfg.Debug)

	errRun := run(cfg, logger)
	if errRun != nil {
		log.Printf("run error: %v", errRun)
		os.Exit(1)
	}
}

func run(cfg *config.Config, loggerй *zap.Logger) error {
	return nil
}

func connectionPostgres() {

}

//go:embed migrations/*.sql
var fsMain embed.FS

func makeMigration(pgConnection string, logger *zap.Logger, f fs.FS) error {
	var d source.Driver
	var errIofs error

	d, errIofs = iofs.New(f, "migrations/")
	if errIofs != nil {
		return fmt.Errorf("error new iofs: %w", errIofs)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, pgConnection)
	if err != nil {
		return fmt.Errorf("error new migrate: %w", err)
	}
	defer m.Close()

	ver, dirty, errGetVersion := m.Version()
	if errGetVersion != nil && !errors.Is(errGetVersion, migrate.ErrNilVersion) {
		return fmt.Errorf("error get version: %w", errGetVersion)
	}
	if dirty {
		m.Force(int(ver - 1))
	}

	errUp := m.Up()
	if errUp != nil {
		if errors.Is(errUp, migrate.ErrNoChange) {
			logger.Info("no change for migrate")
			return nil
		}

		return fmt.Errorf("error up migrate: %w", errUp)
	}

	logger.Info("migrate done")

	return nil
}
