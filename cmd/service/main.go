package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/repository/db"
	"avatar_service/internal/repository/s3storage"
	"avatar_service/internal/service"
	"fmt"
	"log"
	"net"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	S3Client, err := s3storage.New(
		cfg.S3Storage.Endpoint,
		cfg.S3Storage.AccessKeyID,
		cfg.S3Storage.SecretAccessKey,
		true,
	)
	if err != nil {
		log.Fatalln("error s3storage.New: ", err)
	}

	dbRepo, err := db.New(cfg, S3Client)
	if err != nil {
		log.Fatalln("error db.New: ", err)
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
