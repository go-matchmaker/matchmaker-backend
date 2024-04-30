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
	sleepTime = 5 * time.Second
	url       = ""
)

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
	url = getHost

	res := m.Run()

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

func getConnection(ctx context.Context, url string) db.EngineMaker {
	fmt.Println("URL: ", url)
	cfg := setup(url)
	newDB := NewDB(cfg)
	err := newDB.Start(ctx)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	err = newDB.Migration()
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}
	return newDB
}

func stopMigration(db db.EngineMaker) error {
	err := db.Drop()
	if err != nil {
		log.Println("failed to drop db: ", err)
		return err
	}
	return nil
}

func TestConnection(t *testing.T) {
	ctx := context.Background()
	engine := getConnection(ctx, url)

	assert.NotNil(t, engine)
	fmt.Println("Connection established")
}
