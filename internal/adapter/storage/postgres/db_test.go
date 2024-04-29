package psql

import (
	"context"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
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

var (
	sleepTime = time.Millisecond * 500
	conn      = "postgresql"
)

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

func TestMain(m *testing.M) {
	ctx := context.Background()

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
	cfg := setup(getHost)
	newDB := getConnection(ctx, cfg)

	err = newDB.Migration()
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}

	res := m.Run()

	err = newDB.Drop()

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

func TestConnection(t *testing.T) {
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
	if err != nil {
		log.Fatalf("failed to load container: %v", err)
	}
	defer container.Terminate(ctx)

	endpoints, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to get container endpoints: %v", err)
	}
	cfg := setup(endpoints)
	engine := getConnection(ctx, cfg)

	assert.NotNil(t, engine, "Failed to establish connection")
	fmt.Println("Connection established")
}
