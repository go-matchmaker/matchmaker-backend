package dragonfly

import (
	"context"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	sleepTime = 2 * time.Second
	url       = ""
	ctx       = context.Background()
)

func TestMain(m *testing.M) {
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	getHost, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to get endpoint: %s", err)
	}
	url = getHost
	res := m.Run()
	os.Exit(res)
}

func setup(url string) *config.Container {
	endpoint := strings.Split(url, ":")
	host := endpoint[0]
	port, err := strconv.Atoi(endpoint[1])
	address := fmt.Sprintf("%s:%d", host, port)
	if err != nil {
		log.Fatal("failed to convert port to int: ", err)
	}
	return &config.Container{
		Dragonfly: &config.Dragonfly{
			URL: address,
		},
	}
}

func getConnection() cache.EngineMaker {
	cfg := setup(url)
	newCache := NewDragonflyCache(cfg)
	err := newCache.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start cache: %s", err)
	}
	return newCache
}

func TestSetGetDelete(t *testing.T) {
	t.Parallel()

	engine := getConnection()
	time.Sleep(sleepTime)
	require.NotNil(t, engine)
	t.Run("set get delete", func(t *testing.T) {
		key := util.GenerateCacheKey("test", "container")
		value := []byte("test")
		err := engine.Set(ctx, key, value, 0)
		require.NoError(t, err)
		value, err = engine.Get(ctx, key)
		require.NoError(t, err)
		require.Equal(t, "test", string(value))
		err = engine.Delete(ctx, key)
		require.NoError(t, err)
	})
	fmt.Println("Test set get delete done")
}

func TestOthers(t *testing.T) {
	t.Parallel()

	engine := getConnection()
	time.Sleep(sleepTime)
	require.NotNil(t, engine)
	t.Run("set keys delete-prefix", func(t *testing.T) {
		key := util.GenerateCacheKey("test", "container")
		value := []byte("test")
		err := engine.Set(ctx, key, value, 0)
		require.NoError(t, err)
		values, err := engine.Keys(ctx, "*")
		require.NoError(t, err)
		require.NotEmpty(t, values)
		err = engine.DeleteByPrefix(ctx, key)
		require.NoError(t, err)
	})
	fmt.Println("Test set keys delete-prefix done")
}
