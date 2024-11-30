package metrics

import (
	"context"

	"google.golang.org/grpc"
)

type WrappedServerStream struct {
	grpc.ServerStream
	Ctx context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.Ctx
}
