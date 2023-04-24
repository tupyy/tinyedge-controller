package entity

import "time"

type ManifestKind int

func (m ManifestKind) String() string {
	switch m {
	case WorkloadManifestKind:
		return "workload"
	case ConfigurationManifestKind:
		return "configuration"
	default:
		return "unknown"
	}
}

const (
	WorkloadManifestKind = iota
	ConfigurationManifestKind
)

type Manifest interface {
	GetID() string
	GetVersion() string
	GetKind() ManifestKind
	GetHash() string
	GetSelectors() Selectors
	GetNamespaces() []string
	GetSets() []string
	GetDevices() []string
}

type TypeMeta struct {
	// kind to the manifest
	Kind ManifestKind
	// Version
	Version string
}

func (t TypeMeta) GetKind() ManifestKind {
	return t.Kind
}

func (t TypeMeta) GetVersion() string {
	return t.Version
}

type ObjectMeta struct {
	// Id - id of the manifest which is the hash of the filepath
	Id string
	// Labels
	Labels map[string]string
	Hash   string
}

func (o ObjectMeta) GetID() string {
	return o.Id
}

func (o ObjectMeta) GetLabels() map[string]string {
	return o.Labels
}

func (o ObjectMeta) GetHash() string {
	return o.Hash
}

// Workload holds the workload definition.
type Workload struct {
	TypeMeta
	ObjectMeta
	// repository
	Repository Repository
	// path of the manifest file in the local repo
	Path string
	// Description - description of the manifest
	Description string
	// Rootless - set the mode of podman execution: rootless or rootfull
	Rootless bool
	// Secrets - list of secrets without values. Values are retrieve from Vault.
	Secrets []Secret
	// Resources holds the list of file paths
	Resources []string
	// Selectors list of selectors
	Selectors []Selector
	// Devices holds the list of devices' ids which use this manifest
	Devices []string
	// Namespaces hold the list of namespace ids which use this manifest
	Namespaces []string
	// Sets holds the list of sets ids which use this manifest
	Sets []string
}

func (w Workload) GetSelectors() Selectors {
	return w.Selectors
}

func (w Workload) GetNamespaces() []string {
	return w.Namespaces
}

func (w Workload) GetSets() []string {
	return w.Sets
}

func (w Workload) GetDevices() []string {
	return w.Devices
}

// Configuration holds the configuration for a namespace, set or a device.
type Configuration struct {
	TypeMeta
	ObjectMeta
	// repository
	Repository Repository
	// path of the manifest file in the local repo
	Path string
	// list of profiles
	Profiles []Profile
	// HeartbeatPeriod set the heartbeat period of the device
	HeartbeatPeriod time.Duration
	// LogLevel of the device
	LogLevel string
	// Selectors list of selectors
	Selectors []Selector
	// Devices holds the list of devices' ids which use this manifest
	Devices []string
	// Namespaces hold the list of namespace ids which use this manifest
	Namespaces []string
	// Sets holds the list of sets ids which use this manifest
	Sets []string
}

func (c Configuration) GetSelectors() Selectors {
	return c.Selectors
}

func (c Configuration) GetNamespaces() []string {
	return c.Namespaces
}

func (c Configuration) GetSets() []string {
	return c.Sets
}

func (c Configuration) GetDevices() []string {
	return c.Devices
}

type Secret struct {
	Id    string
	Path  string
	Key   string
	Hash  string
	Value string
}

/* DeviceProfile specify all the conditions of a profile:
```yaml
state:
	- perfomance:
		- low: cpu<25%
		- medium: cpu>25%
```
In this example the profile is _perfomance_ and the conditions are _low_ and _medium_.
Each condition's expression is evaluated using Variables.
The expression is only evaluated when all the variables need it by the expression are present in the variable map.
*/
type Profile struct {
	// Name is the name of the profile
	Name string `json:"name"`
	// Conditions holds profile's conditions.
	Conditions []ProfileCondition `json:"conditions"`
}

type ProfileCondition struct {
	// Name is the name of the condition
	Name string `json:"name"`
	// Expression is the boolean expression for the condition
	Expression string `json:"expression"`
}

type SelectorType int

const (
	NamespaceSelector SelectorType = iota
	SetSelector
	DeviceSelector
)

type Selector struct {
	Type  SelectorType
	Value string
}

type Selectors []Selector

func (selectors Selectors) ExtractType(t SelectorType) []string {
	ids := make([]string, 0, len(selectors))
	for _, s := range selectors {
		if s.Type == t {
			ids = append(ids, s.Value)
		}
	}
	return ids
}
