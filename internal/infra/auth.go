package infra

import (
	"context"

	"github.com/s21platform/avatar-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/avatar-service/internal/config"
)

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	_ = info

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata found in context")
	}

	userIDs := md["uuid"]
	if len(userIDs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
	}

	ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])

	return handler(ctx, req)
}

func StreamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "no metadata found in context")
		}

		userIDs := md["uuid"]
		if len(userIDs) == 0 {
			return status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
		}

		ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])

		wrappedStream := &model.ContextServerStream{
			ServerStream: stream,
			Ctx:          ctx,
		}

		return handler(srv, wrappedStream)
	}
}
