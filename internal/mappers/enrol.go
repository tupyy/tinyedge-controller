package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	edgepb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
)

func MapFromEnrolRequest(req *edgepb.EnrolRequest) entity.Device {
	device := entity.Device{
		ID: req.DeviceId,
	}
	if req.Hardware != nil {
		device.HardwareInfo = entity.HardwareInfo{
			Hostname: req.Hardware.GetHostName(),
		}
		if req.Hardware.OsInformation != nil {
			device.HardwareInfo.OsInformation = entity.OsInformation{
				CommitID: req.Hardware.GetOsInformation().GetCommitId(),
			}
		}
		if req.Hardware.SystemVendor != nil {
			device.HardwareInfo.SystemVendor = entity.SystemVendor{
				Manufacturer: req.Hardware.SystemVendor.GetManufacturer(),
				ProductName:  req.Hardware.SystemVendor.GetProductName(),
				SerialNumber: req.Hardware.SystemVendor.GetSerialNumber(),
				Virtual:      req.Hardware.SystemVendor.GetVirtual(),
			}
		}
		device.HardwareInfo.Interfaces = make([]entity.Interface, 0)
		if req.Hardware.Interfaces != nil {
			for _, i := range req.Hardware.GetInterfaces() {
				if i == nil {
					continue
				}
				device.HardwareInfo.Interfaces = append(device.HardwareInfo.Interfaces, entity.Interface{
					Name:        i.GetName(),
					HasCarrier:  i.GetHasCarrier(),
					MacAddress:  i.GetMacAddress(),
					IPV4Address: i.GetIp4Addresses(),
				})
			}
		}
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
