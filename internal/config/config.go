package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres ReadEnvBD
}

type Service struct {
	Port string `env:"POSTGRES_PORT"`
}

type ReadEnvBD struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Println("error cleanenv.ReadEnv: ", err)
	}

	return cfg
}
