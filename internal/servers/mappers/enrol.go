package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	edgepb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
)

func MapFromEnrolRequest(req *edgepb.EnrolRequest) entity.Device {
	device := entity.Device{
		ID: req.DeviceId,
	}
	return device
}

func MapEnrolResponse(enrolStatus entity.EnrolStatus) *edgepb.EnrolResponse {
	resp := &edgepb.EnrolResponse{
		EnrolmentStatus: edgepb.EnrolmentStatus_PENDING,
	}
	switch enrolStatus {
	case entity.EnroledStatus:
		resp.EnrolmentStatus = edgepb.EnrolmentStatus_ENROLED
	case entity.PendingEnrolStatus:
		resp.EnrolmentStatus = edgepb.EnrolmentStatus_PENDING
	case entity.RefusedEnrolStatus:
		resp.EnrolmentStatus = edgepb.EnrolmentStatus_REFUSED
	default:
		resp.EnrolmentStatus = edgepb.EnrolmentStatus_NOT_ENROLED
	}
	return resp
}
