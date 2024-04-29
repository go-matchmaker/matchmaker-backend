package psql

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
	"time"
)

var connURL = ""
var parrallel = false
var sleepTime = time.Millisecond * 500

func migration(uri string) (*migrate.Migrate, error) {
	path := _migrationFilePath

	m, err := migrate.New(path, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate up: %w", err)
	}
	return m, nil
}
func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("foobar"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalln("failed to load container:", err)
	}

	connURL, err = container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalln("failed to get connection string:", err)
	}

	migration, err := migration(connURL)
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}

	res := m.Run()

	migration.Drop()

	os.Exit(res)
}

func getConnection(ctx context.Context, cfg *config.Container) db.EngineMaker {
	newDB := NewDB(cfg)
	err := newDB.Start(ctx)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return newDB
}
