package infra

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
)

func MetricsInterceptor(metrics *pkg.Metrics) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)

		ctx = context.WithValue(ctx, config.KeyMetrics, metrics)
		resp, err := handler(ctx, req)

		if err != nil {
			metrics.Increment(method + "_error")
		}

		metrics.Duration(time.Since(startTime).Milliseconds(), method)

		return resp, err
	}
}

func MetricsStreamInterceptor(metrics *pkg.Metrics) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()

		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)

		wrappedStream := &model.ContextServerStream{
			ServerStream: ss,
			Ctx:          context.WithValue(ss.Context(), config.KeyMetrics, metrics),
		}

		err := handler(srv, wrappedStream)
		if err != nil {
			metrics.Increment(method + "_error")
		}

		metrics.Duration(time.Since(startTime).Milliseconds(), method)

		return err
	}
}
