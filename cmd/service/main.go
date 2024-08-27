package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/repository/db"
	"log"
	"os"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Println("error db.New: ", err)
		os.Exit(1)
	}

	defer dbRepo.Close()
}
