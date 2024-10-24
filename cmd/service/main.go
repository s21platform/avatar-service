package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/repository/db"
	"avatar_service/internal/repository/s3"
	"avatar_service/internal/service"
	"fmt"
	"log"
	"net"

	kafkalib "github.com/s21platform/kafka-lib"

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

	producerNewFriendRegister := kafkalib.NewProducer(cfg.Kafka.Server, cfg.Kafka.AvatarNewSet)

	avatarService := service.New(s3Client, dbRepo, producerNewFriendRegister)
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
