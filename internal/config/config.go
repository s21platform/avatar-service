package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service Service
}

type Service struct {
	Port string `env:"AVATAR_SERVICE_POSTGRES_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Printf("Can not read env variables: %s\n", err)
	}

	return cfg
}
