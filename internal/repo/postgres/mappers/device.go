package mappers

import (
	"database/sql"
	"time"

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

func SetModelToEntity(s []models.SetJoin) entity.Set {
	set := entity.Set{
		Name:        s[0].ID,
		NamespaceID: s[0].NamespaceID,
	}

	idMap := make(uniqueIds)
	ids := make([]string, 0, len(s))
	manifests := make([]string, 0, len(s))
	for _, ss := range s {
		if ss.DeviceId != "" && idMap.exists(ss.DeviceId, "device") {
			ids = append(ids, ss.DeviceId)
			idMap.add(ss.DeviceId, "device")
		}
		if ss.ManifestId != "" && idMap.exists(ss.ManifestId, "manifest") {
			manifests = append(manifests, ss.ManifestId)
			idMap.add(ss.DeviceId, "manifest")
		}
	}

	set.ManifestIDS = manifests
	set.DeviceIDs = ids

	if s[0].DeviceSet.ConfigurationID.Valid {
		set.Configuration = &entity.Configuration{
			ID:              s[0].DeviceSet.ConfigurationID.String,
			HeartbeatPeriod: time.Duration(s[0].ConfigurationHeartbeatPeriodSeconds.Int64 * int64(time.Second)),
			LogLevel:        s[0].ConfigurationLogLevel.String,
		}
	}

	return set
}

func NamespaceModelToEntity(n []models.NamespaceJoin) entity.Namespace {
	namespace := entity.Namespace{
		Name:      n[0].ID,
		IsDefault: false,
	}
	namespace.Configuration = entity.Configuration{
		ID:              n[0].ConfigurationID,
		HeartbeatPeriod: time.Duration(n[0].ConfigurationHeartbeatPeriodSeconds.Int64 * int64(time.Second)),
		LogLevel:        n[0].ConfigurationLogLevel.String,
	}
	if n[0].IsDefault.Valid {
		namespace.IsDefault = n[0].IsDefault.Bool
	}
	idMap := make(uniqueIds)
	sets := make([]string, 0, len(n))
	devices := make([]string, 0, len(n))
	manifests := make([]string, 0, len(n))
	for _, nn := range n {
		if nn.SetId != "" && idMap.exists(nn.SetId, "set") {
			sets = append(sets, nn.SetId)
			idMap.add(nn.SetId, "set")
		}
		if nn.DeviceId != "" && idMap.exists(nn.DeviceId, "device") {
			devices = append(devices, nn.DeviceId)
			idMap.add(nn.DeviceId, "device")
		}
		if nn.ManifestId != "" && idMap.exists(nn.ManifestId, "manifest") {
			manifests = append(manifests, nn.ManifestId)
			idMap.add(nn.ManifestId, "manifest")
		}
	}

	namespace.SetIDs = sets
	namespace.DeviceIDs = devices
	namespace.ManifestIDS = manifests
	return namespace
}
