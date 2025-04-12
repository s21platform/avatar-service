package model

import (
	"context"

	"google.golang.org/grpc"
)

type ContextServerStream struct {
	grpc.ServerStream
	Ctx context.Context
}

func (w *ContextServerStream) Context() context.Context {
	return w.Ctx
}
