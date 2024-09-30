package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/repository/db"
	kafka "avatar_service/internal/repository/kafka/producer/newavatar"
	"avatar_service/internal/repository/s3"
	"avatar_service/internal/service"
	"fmt"
	"log"
	"net"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	s3Client, err := s3.New(cfg)
	if err != nil {
		log.Fatalln("error s3.New: ", err)
	}

	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatalln("error db.New: ", err)
	}
	defer dbRepo.Close()

	kafkaProducer, err := kafka.New(cfg)
	if err != nil {
		log.Println("error kafka.New: ", err)
	}
	defer kafkaProducer.Close()

	avatarService := service.New(s3Client, dbRepo, kafkaProducer)
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
