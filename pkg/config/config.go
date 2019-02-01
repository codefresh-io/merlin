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
			Path string `yaml:"path"`
			Git  struct {
				Owner    string `yaml:"owner"`
				Repo     string `yaml:"repo"`
				Path     string `yaml:"path"`
				Revision string `yaml:"revision"`
			} `yaml:"git"`
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
		Version    string                `yaml:"version"`
		Components []ComponentDescriptor `yaml:"components"`
	}

	ComponentDescriptor struct {
		Name   string   `yaml:"name"`
		Spec   Spec     `yaml:"spec"`
		Values []Values `yaml:values`
	}

	Spec struct {
		Path string `yaml:"path"`
		Git  struct {
			Owner    string `yaml:"owner"`
			Repo     string `yaml:"repo"`
			Path     string `yaml:"path"`
			Revision string `yaml:"revision"`
		} `yaml:"git"`
	}

	Values struct {
		Path string `yaml:"path"`
		Git  struct {
			Owner    string `yaml:"owner"`
			Repo     string `yaml:"repo"`
			Path     string `yaml:"path"`
			Revision string `yaml:"revision"`
		} `yaml:"git"`
	}

	RenderedService struct {
		Commands []Command `yaml:"commands"`
	}

	Command struct {
		Name    string   `yaml:"name"`
		Env     []string `yaml:"env"`
		Program string   `yaml:"program"`
		Args    []string `yaml:"args"`
	}
)

func ReadEnvironmentDescriptor(path string) (*EnvironmentDescriptor, error) {
	cf := &EnvironmentDescriptor{}
	err := utils.ReadFileInto(path, cf)
	return cf, err
}
