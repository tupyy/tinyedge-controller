package mappers

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/manifest"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
)

func DeviceEntityToModel(device entity.Device) models.Device {
	m := models.Device{
		ID:          device.ID,
		NamespaceID: device.NamespaceID,
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

	if device.Configuration != nil {
		m.ConfigurationManifestID = sql.NullString{Valid: true, String: device.Configuration.Id}
	}

	if device.SetID != nil {
		m.DeviceSetID = sql.NullString{Valid: true, String: *device.SetID}
	}

	return m
}

func DeviceToEntity(joins []models.DeviceJoin, readFn manifest.ManifestReader) (entity.Device, error) {
	e := entity.Device{
		ID:          joins[0].ID,
		NamespaceID: joins[0].NamespaceID,
		Registred:   joins[0].Registered,
		EnrolStatus: entity.EnroledStatus.FromString(joins[0].Enroled),
	}
	if joins[0].Registered {
		e.RegisteredAt = joins[0].RegisteredAt
	}

	if e.EnrolStatus == entity.EnroledStatus {
		e.EnroledAt = joins[0].EnroledAt
	}

	if joins[0].CertificateSn.Valid {
		e.CertificateSerialNumber = joins[0].CertificateSn.String
	}
	if joins[0].DeviceSetID.Valid {
		e.SetID = &joins[0].DeviceSetID.String
	}

	idMap := make(uniqueIds)
	manifests := make([]map[string]string, 0, len(joins))
	for _, d := range joins {
		if d.WorkloadID != "" && !idMap.exists(d.WorkloadID, "manifest") {
			manifests = append(manifests, map[string]string{
				"id":   d.WorkloadID,
				"path": d.WorkloadPath,
			})
			idMap.add(d.WorkloadID, "manifest")
		}
	}

	e.Workloads = make([]entity.Workload, 0, len(manifests))
	for _, m := range manifests {
		manifest, err := readManifest(m["path"], m["id"], readFn)
		if err != nil {
			return entity.Device{}, fmt.Errorf("unable to read manifest file %q: %w", m["path"], err)
		}
		e.Workloads = append(e.Workloads, manifest.(entity.Workload))
	}

	return e, nil
}

func DevicesToEntity(joins []models.DeviceJoin, readFn manifest.ManifestReader) ([]entity.Device, error) {
	nmap := make(map[string][]models.DeviceJoin)
	for _, j := range joins {
		_, ok := nmap[j.ID]
		var jj []models.DeviceJoin
		if !ok {
			jj = make([]models.DeviceJoin, 0)
		} else {
			jj = nmap[j.ID]
		}
		jj = append(jj, j)
		nmap[j.ID] = jj
	}

	entities := make([]entity.Device, 0, len(joins))
	for _, v := range nmap {
		d, err := DeviceToEntity(v, readFn)
		if err != nil {
			return entities, err
		}
		entities = append(entities, d)
	}
	return entities, nil
}

func ConfigurationToEntity(c models.Configuration) entity.Configuration {
	e := entity.Configuration{
		ObjectMeta: entity.ObjectMeta{
			Id: c.ID,
		},
		HeartbeatPeriod: time.Duration(c.HeartbeatPeriodSeconds.Int64 * int64(time.Second)),
	}
	if c.LogLevel.Valid {
		e.LogLevel = c.LogLevel.String
	}
	return e
}

func SetsToEntity(sets []models.SetJoin, readFn manifest.ManifestReader) ([]entity.Set, error) {
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
		s, err := SetToEntity(v, readFn)
		if err != nil {
			return entities, err
		}
		entities = append(entities, s)
	}
	return entities, nil
}

func SetToEntity(s []models.SetJoin, readFn manifest.ManifestReader) (entity.Set, error) {
	set := entity.Set{
		Name:        s[0].ID,
		NamespaceID: s[0].NamespaceID,
	}

	if s[0].ConfigurationID != "" {
		// read conf
		c, err := readManifest(s[0].ConfigurationPath, s[0].ConfigurationID, readFn)
		if err != nil {
			return entity.Set{}, err
		}
		conf := c.(entity.Configuration)
		set.Configuration = &conf
	}

	idMap := make(uniqueIds)
	devices := make([]string, 0, len(s))
	manifests := make([]map[string]string, 0, len(s))
	for _, ss := range s {
		if ss.DeviceId != "" && !idMap.exists(ss.DeviceId, "device") {
			devices = append(devices, ss.DeviceId)
			idMap.add(ss.DeviceId, "device")
		}
		if ss.WorkloadID != "" && !idMap.exists(ss.WorkloadID, "manifest") {
			manifests = append(manifests, map[string]string{
				"id":   ss.WorkloadID,
				"path": ss.WorkloadPath,
			})
			idMap.add(ss.WorkloadID, "manifest")
		}
	}

	set.Devices = devices
	set.Workloads = make([]entity.Workload, 0, len(manifests))
	for _, m := range manifests {
		manifest, err := readManifest(m["path"], m["id"], readFn)
		if err != nil {
			return entity.Set{}, fmt.Errorf("unable to read manifest file %q: %w", m["path"], err)
		}
		set.Workloads = append(set.Workloads, manifest.(entity.Workload))
	}

	return set, nil
}

func SetToModel(set entity.Set) models.DeviceSet {
	model := models.DeviceSet{
		ID:          set.Name,
		NamespaceID: set.NamespaceID,
	}
	if set.Configuration != nil {
		model.ConfigurationManifestID = sql.NullString{Valid: true, String: set.Configuration.Id}
	}
	return model
}

func NamespaceToModel(namespace entity.Namespace) models.Namespace {
	model := models.Namespace{
		ID:        namespace.Name,
		IsDefault: sql.NullBool{Valid: true, Bool: namespace.IsDefault},
	}

	if namespace.Configuration != nil {
		model.ConfigurationManifestID = sql.NullString{Valid: true, String: namespace.Configuration.Id}
	}
	return model
}

func NamespacesModelToEntity(namespaces []models.NamespaceJoin, reader manifest.ManifestReader) ([]entity.Namespace, error) {
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
		n, err := NamespaceModelToEntity(v, reader)
		if err != nil {
			return []entity.Namespace{}, err
		}
		entities = append(entities, n)
	}
	return entities, nil
}

func NamespaceModelToEntity(n []models.NamespaceJoin, readFn manifest.ManifestReader) (entity.Namespace, error) {
	namespace := entity.Namespace{
		Name:      n[0].ID,
		IsDefault: false,
	}

	if n[0].ConfigurationPath != "" {
		// read conf
		conf, err := readManifest(n[0].ConfigurationPath, n[0].ConfigurationID, readFn)
		if err != nil {
			return entity.Namespace{}, err
		}

		c := conf.(entity.Configuration)
		namespace.Configuration = &c
	}

	if n[0].Namespace.IsDefault.Valid {
		namespace.IsDefault = n[0].Namespace.IsDefault.Bool
	}

	idMap := make(uniqueIds)
	sets := make([]string, 0, len(n))
	devices := make([]string, 0, len(n))
	manifests := make([]map[string]string, 0, len(n))
	for _, nn := range n {
		if nn.SetId != "" && !idMap.exists(nn.SetId, "set") {
			sets = append(sets, nn.SetId)
			idMap.add(nn.SetId, "set")
		}
		if nn.DeviceId != "" && !idMap.exists(nn.DeviceId, "device") {
			devices = append(devices, nn.DeviceId)
			idMap.add(nn.DeviceId, "device")
		}
		if nn.WorkloadID != "" && !idMap.exists(nn.WorkloadID, "manifest") {
			manifests = append(manifests, map[string]string{
				"id":   nn.WorkloadID,
				"path": nn.WorkloadPath,
			})
			idMap.add(nn.WorkloadID, "manifest")
		}
	}

	namespace.Sets = sets
	namespace.Devices = devices
	namespace.Workloads = make([]entity.Workload, 0, len(manifests))
	for _, m := range manifests {
		manifest, err := readManifest(m["path"], m["id"], readFn)
		if err != nil {
			return entity.Namespace{}, fmt.Errorf("unable to read manifest file %q: %w", m["path"], err)
		}
		namespace.Workloads = append(namespace.Workloads, manifest.(entity.Workload))
	}

	return namespace, nil
}

func readManifest(filepath string, id string, readFn manifest.ManifestReader) (entity.Manifest, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	manifest, err := readFn(bytes.NewBuffer(content), func(m entity.Manifest) entity.Manifest {
		switch v := m.(type) {
		case entity.Workload:
			v.ObjectMeta.Id = id
			return v
		case entity.Configuration:
			v.ObjectMeta.Id = id
			return v
		}
		return m
	})

	if err != nil {
		return nil, err
	}
	return manifest, nil
}
