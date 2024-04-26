package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "config.yml"
)

type (
	Container struct {
		App      `yaml:"app"`
		Log      `yaml:"logger"`
		HTTP     *HTTP     `yaml:"http"`
		Token    *Token    `yaml:"token"`
		Redis    *Redis    `yaml:"redis"`
		PSQL     *PSQL     `yaml:"psql"`
		Settings *Settings `yaml:"settings"`
	}

	App struct {
		Name string `env-required:"true" yaml:"name" env:"APP_NAME"`
	}

	Log struct {
		Level int `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Token struct {
		SymmetricKey string `env-required:"true" yaml:"symmetric_key" env:"TOKEN_SYMMETRIC_KEY"`
		TTL          string `env-required:"true" yaml:"ttl" env:"TOKEN_TTL"`
	}

	Redis struct {
		Host     string `env-required:"true" yaml:"host" env:"REDIS_HOST"`
		Port     int    `env-required:"true" yaml:"port" env:"REDIS_PORT"`
		Password string `env-required:"true" yaml:"password" env:"REDIS_PASSWORD"`
		DbNumber int    `env-required:"true" yaml:"db_number" env:"REDIS_DB_NUMBER"`
	}

	PSQL struct {
		Conn     string `env-required:"true" yaml:"conn" env:"DB_CONN"`
		Host     string `env-required:"true" yaml:"host" env:"DB_HOST"`
		Port     int    `env-required:"true" yaml:"port" env:"DB_PORT"`
		User     string `env-required:"true" yaml:"user" env:"DB_USER"`
		Password string `env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		Name     string `env-required:"true" yaml:"name" env:"DB_NAME"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Kafka struct {
		Host string `env-required:"true" yaml:"host" env:"KAFKA_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"KAFKA_PORT"`
	}

	Settings struct {
		ServerReadTimeout int `env-required:"true" yaml:"server_read_timeout" env:"SERVER_READ_TIMEOUT"`
		PSQLConnAttempts  int `env-required:"true" yaml:"psql_conn_attempts" env:"PSQL_CONN_ATTEMPTS"`
		PSQLConnTimeout   int `env-required:"true" yaml:"psql_conn_timeout" env:"PSQL_CONN_TIMEOUT"`
	}
)

func NewConfig() (*Container, error) {
	cfg := new(Container)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
