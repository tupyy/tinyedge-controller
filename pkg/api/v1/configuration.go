package v1

type Configuration struct {
	Kind        string `yaml:"kind" validate:"required"`
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
