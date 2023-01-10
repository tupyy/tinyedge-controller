package entity

type RelationType int

const (
	NamespaceRelationType RelationType = iota
	SetRelationType
	DeviceRelationType
)

type ReferenceRelation struct {
	Type       RelationType
	ResourceID string
	ManifestID string
}

func NewNamespaceRelation(namespaceId, manifestId string) ReferenceRelation {
	return ReferenceRelation{
		Type:       NamespaceRelationType,
		ResourceID: namespaceId,
		ManifestID: manifestId,
	}
}

func NewSetRelation(setId, manifestId string) ReferenceRelation {
	return ReferenceRelation{
		Type:       SetRelationType,
		ResourceID: setId,
		ManifestID: manifestId,
	}
}

func NewDeviceRelation(deviceId, manifestId string) ReferenceRelation {
	return ReferenceRelation{
		Type:       DeviceRelationType,
		ResourceID: deviceId,
		ManifestID: manifestId,
	}
}
