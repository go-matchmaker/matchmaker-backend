package psql

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_migrationFilePath = "internal/adapter/storage/postgres/migration"
	_defaultAttempts   = 5
	_defaultTimeout    = time.Second
)

func (ps *postgres) Migration() {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)
	zap.S().Infoln("Migrate: up started")
	connUrl := ps.GetURL()
	for attempts > 0 {
		m, err = migrate.New(fmt.Sprintf("file://%s", _migrationFilePath), connUrl)
		if err == nil || errors.Is(err, migrate.ErrNoChange) {
			break
		}

		zap.S().Infoln("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		zap.S().Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.S().Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		zap.S().Infoln("Migrate: no change")
	}

	zap.S().Infoln("Migrate: up finished")
}
