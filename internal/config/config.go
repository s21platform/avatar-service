package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service   Service
	Postgres  Postgres
	S3Storage S3Storage
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

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalln("error cleanenv.ReadEnv: ", err)
	}

	return cfg
}
