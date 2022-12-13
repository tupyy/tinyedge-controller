package servers

import (
	"context"
	"errors"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/mappers"
	errService "github.com/tupyy/tinyedge-controller/internal/services/common"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	pb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EdgeServer struct {
	pb.UnimplementedEdgeServiceServer
	edgeService *edge.Service
}

func (e *EdgeServer) New(edgeService *edge.Service) *EdgeServer {
	return &EdgeServer{edgeService: edgeService}
}

func (e *EdgeServer) Enrol(ctx context.Context, req *pb.EnrolRequest) (*pb.EnrolResponse, error) {
	newDevice := mappers.MapFromEnrolRequest(req)
	enrolStatus, err := e.edgeService.Enrol(ctx, newDevice)
	if err != nil {
		return &pb.EnrolResponse{
			EnrolmentStatus: pb.EnrolmentStatus_REFUSED,
		}, status.Error(codes.Internal, "internal error")
	}

	return mappers.MapEnrolResponse(enrolStatus), nil
}

func (e *EdgeServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	// check if it is enroled
	IsEnroled, err := e.edgeService.IsEnroled(ctx, req.DeviceId)
	if err != nil {
		if errors.Is(err, errService.ErrDeviceNotFound) {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("device %q not found. Please enrol first", req.GetDeviceId()))
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	if !IsEnroled {
		return nil, status.Errorf(codes.InvalidArgument, "device %q is not enroled yet. Please enrol the device.", req.DeviceId)
	}

	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (e *EdgeServer) GetConfiguration(ctx context.Context, req *pb.ConfigurationRequest) (*pb.ConfigurationResponse, error) {
	// check if it is registered
	isRegistered, err := e.edgeService.IsRegistered(ctx, req.DeviceId)
	if err != nil {
		if errors.Is(err, errService.ErrDeviceNotFound) {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("device %q not found", req.GetDeviceId()))
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	if !isRegistered {
		return nil, status.Errorf(codes.InvalidArgument, "device is not registered")
	}

	configuration, err := e.edgeService.GetConfiguration(ctx, req.DeviceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return mappers.MapConfigurationToProto(configuration), nil
}

func (e *EdgeServer) Heartbeat(ctx context.Context, req *common.HeartbeatInfo) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
