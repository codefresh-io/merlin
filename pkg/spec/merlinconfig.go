package spec

type (
	// Config representation of config.yaml file of one environment
	Config struct {
		Name    string            `yaml:"name"`
		Shell   string            `yaml:"shell"`
		Values  map[string]string `yaml:"value"`
		Cluster Cluster           `yaml:"cluster"`
	}
	Cluster struct {
		Path      string `yaml:"path"`
		Context   string `yaml:"context"`
		Namespace string `yaml:"namespace"`
	}
)
