package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	kafkalib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/infra"
	"github.com/s21platform/avatar-service/internal/repository/postgres"
	"github.com/s21platform/avatar-service/internal/repository/s3"
	"github.com/s21platform/avatar-service/internal/service"
	"github.com/s21platform/avatar-service/pkg/avatar"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	s3Client := s3.New(cfg)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "avatar", cfg.Platform.Env)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create metrics object: %v", err))
		log.Fatal("failed to create metrics object: ", err)
	}

	producerNewUserAvatarRegister := kafkalib.NewProducer(cfg.Kafka.Server, cfg.Kafka.UserNewSet)
	producerNewSocietyAvatarRegister := kafkalib.NewProducer(cfg.Kafka.Server, cfg.Kafka.SocietyNewSet)

	avatarService := service.New(s3Client, dbRepo, producerNewUserAvatarRegister, producerNewSocietyAvatarRegister, cfg.S3Storage.BucketName)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
			infra.Logger(logger),
		),
		grpc.ChainStreamInterceptor(
			infra.MetricsStreamInterceptor(metrics),
			infra.StreamLogger(logger),
		),
	)

	avatar.RegisterAvatarServiceServer(grpcServer, avatarService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to start TCP listener: %v", err))
	}

	if err = grpcServer.Serve(listener); err != nil {
		logger.Error(fmt.Sprintf("failed to start gRPC listener: %v", err))
	}
}
