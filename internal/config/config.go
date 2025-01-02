package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service   Service
	Postgres  Postgres
	S3Storage S3Storage
	Kafka     Kafka
	Metrics   Metrics
	Platform  Platform
}

type Service struct {
	Port string `env:"AVATAR_SERVICE_PORT"`
}

type Postgres struct {
	User     string `env:"AVATAR_SERVICE_POSTGRES_USER"`
	Password string `env:"AVATAR_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"AVATAR_SERVICE_POSTGRES_DB"`
	Host     string `env:"AVATAR_SERVICE_POSTGRES_HOST"`
	Port     string `env:"AVATAR_SERVICE_POSTGRES_PORT"`
}

type S3Storage struct {
	Endpoint        string `env:"AVATAR_SERVICE_S3_STORAGE_ENDPOINT"`
	AccessKeyID     string `env:"AVATAR_SERVICE_S3_STORAGE_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AVATAR_SERVICE_S3_STORAGE_ACCESS_KEY_SECRET"`
}

type Kafka struct {
	AvatarNewSet string `env:"AVATAR_SET_NEW"`
	Server       string `env:"KAFKA_SERVER"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}

	return cfg
}
