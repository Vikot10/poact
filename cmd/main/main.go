package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Vikot10/viarticles/internal/config"
)

var version = "undefined"

func main() {
	cfg := config.MustLoad()

	logger := mustCreateLogger(cfg.Debug)

	errRun := run(cfg, logger)
	if errRun != nil {
		log.Printf("run error: %v", errRun)
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger.Info("start", zap.String("version", version))

	ln, errLn := net.Listen("tcp", cfg.Address)
	if errLn != nil {
		return fmt.Errorf("listen error: %w", errLn)
	}
	defer ln.Close()

	dbPool, errDB := connectionPostgres(ctx, cfg.Postgres)
	if errDB != nil {
		return fmt.Errorf("db error: %w", errDB)
	}
	defer dbPool.Close()

	if cfg.Postgres.NeedMigrate {
		err := makeMigration(buildPGConnectionString(cfg.Postgres), logger, fsMain)
		if err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
	}
	return nil
}

func buildPGConnectionString(cfgPg config.Postgres) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
		cfgPg.Username,
		cfgPg.Password,
		cfgPg.Host,
		cfgPg.Port,
		cfgPg.Database,
		cfgPg.SSLMode,
		cfgPg.SSLCertPath)
}

func connectionPostgres(ctx context.Context, cfgPg config.Postgres) (*pgxpool.Pool, error) {
	stringConnection := buildPGConnectionString(cfgPg)

	dbPool, errDB := pgxpool.New(ctx, stringConnection)
	if errDB != nil {
		return nil, fmt.Errorf("error create pg pool: %w", errDB)
	}

	errDB = dbPool.Ping(ctx)
	if errDB != nil {
		return nil, fmt.Errorf("error ping pg pool: %w", errDB)
	}

	return dbPool, nil
}

//go:embed migrations/*.sql
var fsMain embed.FS

func makeMigration(pgConnection string, logger *zap.Logger, f fs.FS) error {
	var d source.Driver
	var errIofs error

	d, errIofs = iofs.New(f, "migrations")
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
