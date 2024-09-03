package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/repository/db"
	"avatar_service/internal/service"
	"fmt"
	"log"
	"net"
	"os"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Println("error db.New: ", err)
		os.Exit(1)
	}

	defer dbRepo.Close()

	avatarService := service.New(dbRepo)
	grpcServer := grpc.NewServer()
	avatarproto.RegisterAvatarServiceServer(grpcServer, avatarService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Println("error net.Listen: ", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Println("error grpcServer.Serve: ", err)
	}

}
