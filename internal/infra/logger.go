package infra

import (
	"context"

	"google.golang.org/grpc"

	logger "github.com/s21platform/logger-lib"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
)

func Logger(logger *logger.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, config.KeyLogger, logger)
		return handler(ctx, req)
	}
}

func StreamLogger(logger *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		ctx = context.WithValue(ctx, config.KeyLogger, logger)
		return handler(srv, &model.ContextServerStream{ServerStream: stream, Ctx: ctx})
	}
}
