package entity

type RelationType int

const (
	NamespaceRelationType RelationType = iota
	SetRelationType
	DeviceRelationType
)

type Relation struct {
	Type       RelationType
	ResourceID string
	ManifestID string
}

func NewNamespaceRelation(namespaceId, manifestId string) Relation {
	return Relation{
		Type:       NamespaceRelationType,
		ResourceID: namespaceId,
		ManifestID: manifestId,
	}
}

func NewSetRelation(setId, manifestId string) Relation {
	return Relation{
		Type:       SetRelationType,
		ResourceID: setId,
		ManifestID: manifestId,
	}
}

func NewDeviceRelation(deviceId, manifestId string) Relation {
	return Relation{
		Type:       DeviceRelationType,
		ResourceID: deviceId,
		ManifestID: manifestId,
	}
}
