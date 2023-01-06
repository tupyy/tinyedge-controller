package entity

import (
	"time"
)

type ConfigurationResponse struct {
	Configuration Configuration
	Workloads     []Workload
	Secrets       []Secret
}

func (c ConfigurationResponse) Hash() string {
	return "hash"
}

type Configuration struct {
	ID string
	// list of profiles
	Profiles []Profile
	// HeartbeatConfiguration hold the configuration of hearbeat
	HeartbeatPeriod time.Duration
	LogLevel        string
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
