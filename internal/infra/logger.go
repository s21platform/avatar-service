package infra

import (
	"avatar_service/internal/config"
	"avatar_service/internal/model"
	"context"

	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc"
)

func Logger(logger *logger_lib.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, config.KeyLogger, logger)
		return handler(ctx, req)
	}
}

func StreamLogger(logger *logger_lib.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()

		ctx = context.WithValue(ctx, config.KeyLogger, logger)

		return handler(srv, &model.WrappedServerStream{ServerStream: stream, Ctx: ctx})
	}
}
