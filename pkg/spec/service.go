package spec

type (
	Service struct {
		Name        string   `yaml:"name"`
		Ports       []Port   `yaml:"ports"`
		Environment []EnvVar `yaml:"environment"`
	}

	Port struct {
		Name    string `yaml:"name"`
		EnvVar  string `yaml:"envVar"`
		Default int    `yaml:"default"`
	}

	EnvVar struct {
		Name    string `yaml:"name"`
		Default string `yaml:"default"`
	}
)
