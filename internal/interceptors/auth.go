package interceptors

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/services/auth"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	deviceIDKey = "device_id"
)

func AuthInterceptor(auth *auth.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		peer, _ := peer.FromContext(ctx)
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)

		if len(tlsInfo.State.PeerCertificates) == 0 {
			return common.Empty{}, status.Errorf(codes.PermissionDenied, "missing peer certificates")
		}

		deviceID, err := getDeviceIDFromContext(ctx)
		if err != nil {
			return common.Empty{}, err
		}

		newCtx, err := auth.Auth(ctx, info.FullMethod, deviceID, tlsInfo.State.PeerCertificates)
		if err != nil {
			zap.S().Errorf("unable to authenticate device", "error", err)
			return common.Empty{}, status.Errorf(codes.PermissionDenied, err.Error())
		}

		return handler(newCtx, req)
	}
}

func getDeviceIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.PermissionDenied, "cannot retrieve metadata from context")
	}

	data := md.Get(deviceIDKey)
	if len(data) == 0 {
		return "", status.Errorf(codes.PermissionDenied, "cannot retrieve device id from metadata")
	}

	if data[0] == "" {
		return "", status.Errorf(codes.PermissionDenied, "cannot retrieve device id from metadata")
	}

	return data[0], nil
}
