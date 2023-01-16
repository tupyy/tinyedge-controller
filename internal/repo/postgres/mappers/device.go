package mappers

import (
	"database/sql"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
)

func MapEntityToModel(device entity.Device) models.Device {
	m := models.Device{
		ID:          device.ID,
		NamespaceID: sql.NullString{Valid: true, String: device.NamespaceID},
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
		NamespaceID: device.NamespaceID.String,
		Registred:   device.Registered,
		EnroledAt:   device.EnroledAt,
		EnrolStatus: entity.EnroledStatus.FromString(device.Enroled),
	}
	if device.CertificateSn.Valid {
		e.CertificateSerialNumber = device.CertificateSn.String
	}
	if device.DeviceSetID.Valid {
		e.SetID = &device.DeviceSetID.String
	}
	return e
}

func SetsToEntity(sets []models.SetJoin) []entity.Set {
	nmap := make(map[string][]models.SetJoin)
	for _, n := range sets {
		_, ok := nmap[n.ID]
		var nn []models.SetJoin
		if !ok {
			nn = make([]models.SetJoin, 0)
		} else {
			nn = nmap[n.ID]
		}
		nn = append(nn, n)
		nmap[n.ID] = nn
	}

	entities := make([]entity.Set, 0, len(sets))
	for _, v := range nmap {
		entities = append(entities, SetToEntity(v))
	}
	return entities
}

func SetToEntity(s []models.SetJoin) entity.Set {
	set := entity.Set{
		Name:        s[0].ID,
		NamespaceID: s[0].NamespaceID,
	}

	idMap := make(uniqueIds)
	ids := make([]string, 0, len(s))
	manifests := make([]string, 0, len(s))
	for _, ss := range s {
		if ss.DeviceId != "" && !idMap.exists(ss.DeviceId, "device") {
			ids = append(ids, ss.DeviceId)
			idMap.add(ss.DeviceId, "device")
		}
		if ss.ManifestId != "" && !idMap.exists(ss.ManifestId, "manifest") {
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

func NamespacesModelToEntity(namespaces []models.NamespaceJoin) []entity.Namespace {
	nmap := make(map[string][]models.NamespaceJoin)
	for _, n := range namespaces {
		_, ok := nmap[n.ID]
		var nn []models.NamespaceJoin
		if !ok {
			nn = make([]models.NamespaceJoin, 0)
		} else {
			nn = nmap[n.ID]
		}
		nn = append(nn, n)
		nmap[n.ID] = nn
	}

	entities := make([]entity.Namespace, 0, len(namespaces))
	for _, v := range nmap {
		entities = append(entities, NamespaceModelToEntity(v))
	}
	return entities
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
	if n[0].Namespace.IsDefault.Valid {
		namespace.IsDefault = n[0].Namespace.IsDefault.Bool
	}
	idMap := make(uniqueIds)
	sets := make([]string, 0, len(n))
	devices := make([]string, 0, len(n))
	manifests := make([]string, 0, len(n))
	for _, nn := range n {
		if nn.SetId != "" && !idMap.exists(nn.SetId, "set") {
			sets = append(sets, nn.SetId)
			idMap.add(nn.SetId, "set")
		}
		if nn.DeviceId != "" && !idMap.exists(nn.DeviceId, "device") {
			devices = append(devices, nn.DeviceId)
			idMap.add(nn.DeviceId, "device")
		}
		if nn.ManifestId != "" && !idMap.exists(nn.ManifestId, "manifest") {
			manifests = append(manifests, nn.ManifestId)
			idMap.add(nn.ManifestId, "manifest")
		}
	}

	namespace.SetIDs = sets
	namespace.DeviceIDs = devices
	namespace.ManifestIDS = manifests
	return namespace
}
