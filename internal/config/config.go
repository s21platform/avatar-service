package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
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

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalln("error cleanenv.ReadEnv: ", err)
	}

	return cfg
}
