package interceptors

import (
	"context"
	"regexp"

	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	pb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IDValidationInterceptor returns an interceptor which will validate the id of the device.
func IDValidationInterceptor(validationFn func(id string) bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		returnFn := func(resp interface{}) (interface{}, error) {
			return resp, status.Error(codes.InvalidArgument, "device id is invalid")
		}

		switch v := req.(type) {
		case pb.EnrolRequest:
			if !validationFn(v.DeviceId) {
				return returnFn(&pb.EnrolResponse{EnrolmentStatus: common.EnrolmentStatus_REFUSED})
			}
		case pb.ConfigurationRequest:
			if !validationFn(v.DeviceId) {
				return returnFn(&pb.ConfigurationResponse{})
			}
		case pb.RegistrationRequest:
			if !validationFn(v.DeviceId) {
				return returnFn(&pb.RegistrationResponse{})
			}
		case common.HeartbeatInfo:
			if !validationFn(v.DeviceId) {
				return returnFn(&common.Empty{})
			}
		}

		return handler(ctx, req)
	}
}

func ValidateDeviceID(id string) bool {
	validID := regexp.MustCompile(`^[^±!@£$%&*+§¡€#¢§¶•ªº«\\/<>?:;|=,\s]*$`)
	if id == "" {
		return false
	}
	if len(id) > 256 || !validID.MatchString(id) {
		return false
	}
	return true
}
