package entity

import "time"

const (
	// name of the manifest work
	ManifestWorkFilename = "manifest_work.yaml"
)

const (
	WorkloadManifestKind      = "workload"
	ConfigurationManifestKind = "configuration"
)

type Manifest interface {
	GetID() string
	GetVersion() string
	GetName() string
	GetKind() string
	GetPath() string
	GetRepository() Repository
}

type TypeMeta struct {
	// kind to the manifest
	Kind string
	// Version
	Version string
}

type ObjectMeta struct {
	// Name - name of the manifest
	Name string
	// Labels
	Labels map[string]string
	// Id - id of the manifest which is the hash of the filepath
	Id string
}

// Workload holds the workload definition.
type Workload struct {
	TypeMeta
	ObjectMeta
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
	// repository
	Repository Repository
}

func (w Workload) GetName() string {
	return w.ObjectMeta.Name
}

func (w Workload) GetKind() string {
	return w.TypeMeta.Kind
}

func (w Workload) GetID() string {
	return w.ObjectMeta.Id
}

func (w Workload) GetVersion() string {
	return w.TypeMeta.Version
}

func (w Workload) GetPath() string {
	return w.Path
}

func (w Workload) GetRepository() Repository {
	return w.Repository
}

type Secret struct {
	Id    string
	Path  string
	Key   string
	Hash  string
	Value string
}

// Configuration holds the configuration for a namespace, set or a device.
type Configuration struct {
	TypeMeta
	ObjectMeta
	// list of profiles
	Profiles []Profile
	// HeartbeatPeriod set the heartbeat period of the device
	HeartbeatPeriod time.Duration
	// LogLevel of the device
	LogLevel string
	// Selectors list of selectors
	Selectors []Selector
}

func (c Configuration) GetName() string {
	return c.ObjectMeta.Name
}

func (c Configuration) GetKind() string {
	return c.TypeMeta.Kind
}

func (c Configuration) GetID() string {
	return c.ObjectMeta.Id
}

func (c Configuration) GetVersion() string {
	return c.TypeMeta.Version
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
