package mappers

import (
	"database/sql"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"
)

func MapEntityToModel(device entity.Device) models.Device {
	m := models.Device{
		ID:          device.ID,
		NamespaceID: sql.NullString{Valid: true, String: device.Namespace},
		Registered:  device.Registred,
		Enroled:     device.EnrolStatus.String(),
		EnroledAt:   device.EnroledAt,
	}

	if device.EnrolStatus == entity.EnroledStatus {
		m.EnroledAt = device.EnroledAt
	}

	if device.Registred {
		m.RegisteredAt = device.RegisteredAt
		if device.CertificateSerialNumber != "" {
			m.CertificateSn = sql.NullString{Valid: true, String: device.CertificateSerialNumber}
		}
	}

	return m
}

func MapModelToEntity(device models.Device) entity.Device {
	e := entity.Device{
		ID:          device.ID,
		Namespace:   device.NamespaceID.String,
		Registred:   device.Registered,
		EnroledAt:   device.EnroledAt,
		EnrolStatus: entity.EnroledStatus.FromString(device.Enroled),
	}
	if device.CertificateSn.Valid {
		e.CertificateSerialNumber = device.CertificateSn.String
	}
	return e
}
