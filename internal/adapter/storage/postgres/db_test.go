package psql

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"strconv"
	"strings"

	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

const ()

var (
	sleepTime = 2 * time.Second
	url       = ""
	ctx       = context.Background()
)

func TestMain(m *testing.M) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalln("failed to load container:", err)
	}
	getHost, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalln("failed to load container:", err)
	}
	url = getHost
	migrateURL := "postgres://testuser:testpassword@" + url + "/testdb?sslmode=disable"
	fmt.Println("url: ", migrateURL, _migrationFilePath)
	migrateUp, err := MigrateUp(migrateURL)
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}
	res := m.Run()
	migrateUp.Drop()

	os.Exit(res)
}

func setup(url string) *config.Container {
	endpoint := strings.Split(url, ":")
	host := endpoint[0]
	port, err := strconv.Atoi(endpoint[1])
	if err != nil {
		log.Fatal("failed to convert port to int: ", err)
	}
	return &config.Container{
		PSQL: &config.PSQL{
			Conn:     "postgresql",
			Host:     host,
			Port:     port,
			User:     "testuser",
			Password: "testpassword",
			Name:     "testdb",
		},
		Settings: &config.Settings{
			PSQLConnTimeout:  5,
			PSQLConnAttempts: 5},
	}
}

func getConnection() db.EngineMaker {
	cfg := setup(url)
	newDB := NewDB(cfg)
	err := newDB.Start(ctx)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return newDB
}

func cleanUp() {
	getConnection().Execute(ctx, "DELETE FROM users")
}

func MigrateUp(uri string) (*migrate.Migrate, error) {
	source, err := iofs.New(migrationsFS, _migrationFilePath)
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithSourceInstance("iofs", source, uri)
	if err != nil {

		return nil, err
	}
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate up: %w", err)
	}
	return migration, nil
}

func TestConnection(t *testing.T) {
	engine := getConnection()

	assert.NotNil(t, engine)
	fmt.Println("Connection established")
}
