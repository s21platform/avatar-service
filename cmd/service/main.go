package main

import (
	"avatar_service/internal/config"
	"avatar_service/internal/infra"
	"avatar_service/internal/repository/db"
	"avatar_service/internal/repository/s3"
	"avatar_service/internal/service"
	"fmt"
	"log"
	"net"

	kafkalib "github.com/s21platform/kafka-lib"
	"github.com/s21platform/metrics-lib/pkg"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	s3Client, err := s3.New(cfg)
	if err != nil {
		log.Fatalln("failed to create S3 client: ", err)
	}

	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatalln("failed to initialize database repository: ", err)
	}
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "avatar", cfg.Platform.Env)
	if err != nil {
		log.Fatalln("failed to create metrics object: ", err)
	}

	producerNewFriendRegister := kafkalib.NewProducer(cfg.Kafka.Server, cfg.Kafka.AvatarNewSet)

	avatarService := service.New(s3Client, dbRepo, producerNewFriendRegister)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
		),
		grpc.ChainStreamInterceptor(
			infra.MetricsStreamInterceptor(metrics),
		),
	)

	avatarproto.RegisterAvatarServiceServer(grpcServer, avatarService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Println("failed to start TCP listener: ", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Println("failed to start grpc server: ", err)
	}
}
