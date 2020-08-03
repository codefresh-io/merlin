package spec

const (
	// PortGenerationStrategyRandom generate random port from connecting to the
	// service or debugging
	PortGenerationStrategyRandom = "random"
)

type (
	// Config representation of config.yaml file of one environment
	Config struct {
		Name                   string                 `yaml:"name"`
		Shell                  string                 `yaml:"shell"`
		Values                 map[string]string      `yaml:"value"`
		Cluster                Cluster                `yaml:"cluster"`
		PortGenerationStrategy PortGenerationStrategy `yaml:"portGenerationStrategy"`
	}
	Cluster struct {
		Path      string `yaml:"path"`
		Context   string `yaml:"context"`
		Namespace string `yaml:"namespace"`
	}

	PortGenerationStrategy struct {
		PortForward *string `yaml:"portForward,omitempty"`
		Debug       *string `yaml:"debug,omitempty"`
	}
)
