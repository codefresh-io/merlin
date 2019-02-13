package config

import (
	"github.com/codefresh-io/merlin/pkg/utils"
)

const (
	MerlinFileName      = "merlin"
	MerlinFileExtention = "yaml"
	MerlinFilePath      = ".codefresh/.dev"
)

type (
	Config struct {
		Name      string `yaml:"name"`
		Codefresh struct {
			Context string `yaml:"context"`
			Path    string `yaml:"path"`
		} `yaml:"codefresh"`
		Environment struct {
			Spec   Spec     `yaml:"spec"`
			Values []Values `yaml:"values"`
		} `yaml:"environment"`
		Kube struct {
			Context   string `yaml:"context"`
			Path      string `yaml:"path"`
			Namespace string `yaml:"namespace"`
		} `yaml:"kube"`
		Github struct {
			Token       string `yaml:"token"`
			PemFilePath string `yaml:"pemFilePath"`
		} `yaml:"github"`
		Cache struct {
			Path string `yaml:"path"`
		} `yaml:"cache"`
	}

	EnvironmentDescriptor struct {
		Version    string      `yaml:"version"`
		Components []Component `yaml:"components"`
		Operators  []Operator  `yaml:"operators"`
	}

	Component struct {
		Name   string   `yaml:"name"`
		Spec   Spec     `yaml:"spec"`
		Values []Values `yaml:values`
	}

	Spec struct {
		Path string `yaml:"path"`
		Git  Git    `yaml:"git"`
	}

	Values struct {
		Path string `yaml:"path"`
		Git  Git    `yaml:"git"`
	}

	Git struct {
		Owner    string `yaml:"owner"`
		Repo     string `yaml:"repo"`
		Path     string `yaml:"path"`
		Revision string `yaml:"revision"`
	}

	ComponentDescriptor struct {
		Operators []Operator `yaml:"operators"`
	}

	Operator struct {
		Name        string `yaml:"name"`
		Type        string `yaml:"type"`
		Description string `yaml:"description"`
		Spec        struct {
			Pipeline  string   `yaml:"pipeline"`
			Branch    string   `yaml:"branch"`
			Variables []string `yaml:"variables"`

			Env     []string `yaml:"env"`
			Program string   `yaml:"program"`
			Args    []string `yaml:"args"`

			Detached bool `yaml:"detached"`
		} `yaml:"spec"`
	}
)

func ReadEnvironmentDescriptor(path string) (*EnvironmentDescriptor, error) {
	cf := &EnvironmentDescriptor{}
	err := utils.ReadFileInto(path, cf)
	return cf, err
}
