package main

import (
	"avatar_service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg
}
