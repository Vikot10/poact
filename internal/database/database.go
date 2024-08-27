package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed database/*.sql
var fsMain embed.FS

func MakeMigration(pgConnection string, logger *zap.Logger) error {
	var d source.Driver
	var errIofs error

	d, errIofs = iofs.New(fsMain, "database")
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
