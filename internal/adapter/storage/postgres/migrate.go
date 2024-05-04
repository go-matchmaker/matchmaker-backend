package psql

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migration/*.sql
var migrationsFS embed.FS

const (
	_migrationFilePath = "migration"
	_defaultAttempts   = 5
	_defaultTimeout    = time.Second
)

func (ps *pdb) Migration() error {
	var (
		attempts = _defaultAttempts
		err      error
	)
	zap.S().Infoln("Migrate: up started")
	for attempts > 0 {
		err = ps.migrationSettings()
		if err == nil {
			break
		}
		zap.S().Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		zap.S().Fatal("Migrate error: %s", err)
	}
	return nil
}

func (ps *pdb) migrationSettings() error {
	connURL := ps.getURL()
	fmt.Println("URL:", connURL)
	source, err := iofs.New(migrationsFS, _migrationFilePath)
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithSourceInstance("iofs", source, connURL)
	if err != nil {
		return err
	}
	err = migration.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			zap.S().Info("Migrate: no changes")
			return nil
		}
	}
	return nil
}

func (ps *pdb) Drop() error {
	connURL := ps.getURL()
	source, err := iofs.New(migrationsFS, _migrationFilePath)
	if err != nil {
		zap.S().Fatal("Migrate: drop error: %s", err)
	}
	migration, err := migrate.NewWithSourceInstance("iofs", source, connURL)
	if err != nil {
		zap.S().Fatal("Migrate: drop error: %s", err)
	}
	err = migration.Down()
	if err != nil {
		zap.S().Fatal("Migrate: drop error: %s", err)
	}
	zap.S().Info("Migrate: drop completed")
	return nil
}
