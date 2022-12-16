package servers

import (
	"context"
	"errors"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/mappers"
	errService "github.com/tupyy/tinyedge-controller/internal/services/common"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	pb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EdgeServer struct {
	pb.UnimplementedEdgeServiceServer
	edgeService *edge.Service
}

func NewEdgeServer(edgeService *edge.Service) *EdgeServer {
	return &EdgeServer{edgeService: edgeService}
}

func (e *EdgeServer) Enrol(ctx context.Context, req *pb.EnrolRequest) (*pb.EnrolResponse, error) {
	enrolStatus, err := e.edgeService.Enrol(ctx, req.DeviceId)
	if err != nil {
		zap.S().Errorw("unable to enrol device", "error", err, "device_id", req.DeviceId)
		return &pb.EnrolResponse{
			EnrolmentStatus: pb.EnrolmentStatus_REFUSED,
		}, status.Error(codes.Internal, "internal error")
	}

	if enrolStatus == entity.RefusedEnrolStatus {
		return mappers.MapEnrolResponse(enrolStatus), status.Errorf(codes.PermissionDenied, "device %q enrol request has been denied", req.DeviceId)
	}

	return mappers.MapEnrolResponse(enrolStatus), nil
}

func (e *EdgeServer) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	certificate, err := e.edgeService.Register(ctx, req.DeviceId, req.CertificateRequest)
	if err != nil {
		if errors.Is(err, errService.ErrDeviceNotEnroled) {
			return nil, status.Errorf(codes.InvalidArgument, "unable to register device %q. Device is not enroled", req.DeviceId)
		}

		if errors.Is(err, errService.ErrDeviceNotFound) {
			return nil, status.Errorf(codes.NotFound, "device %q not found. Please enrol the device first.", req.DeviceId)
		}

		zap.S().Errorw("unable to register device", "error", err, "device_id", req.DeviceId)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.RegistrationResponse{Certificate: string(certificate.CertificatePEM)}, nil
}

func (e *EdgeServer) GetConfiguration(ctx context.Context, req *pb.ConfigurationRequest) (*pb.ConfigurationResponse, error) {
	// guarded by the real device certificate
	configuration, err := e.edgeService.GetConfiguration(ctx, req.DeviceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return mappers.MapConfigurationToProto(configuration), nil
}

func (e *EdgeServer) Heartbeat(ctx context.Context, req *common.HeartbeatInfo) (*common.Empty, error) {
	return &common.Empty{}, nil
}
