package v1

type Manifest struct {
	Version     string     `yaml:"version"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Selector    Selector   `yaml:"selectors"`
	Secrets     []Secret   `yaml:"secrets"`
	Resources   []Resource `yaml:"resources"`
}

type Selector struct {
	Namespaces []string `yaml:"namespaces"`
	Sets       []string `yaml:"sets"`
	Devices    []string `yaml:"devices"`
}

type Secret struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Key  string `yaml:"key"`
}

type Resource struct {
	Ref string `yaml:"$ref"`
}
