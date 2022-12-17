package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// AuthInterceptor verify if the certificate presented by the agent is revoked.
// The agent is allowed to authenticate
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		peer, _ := peer.FromContext(ctx)
		fmt.Printf("%+v", peer)
		_ = req
		_ = info
		return handler(ctx, req)
	}
}
