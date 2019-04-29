package spec

import (
	"encoding/json"
	"fmt"
)

type (
	ActiveConfig struct {
		Name          string `json:"name"`
		EnvironmentJS string `json:"environment"`
		Codefresh     struct {
			Path    string `json:"path"`
			Context string `json:"context"`
		} `json:"codefresh"`
		Kubernetes struct {
			Path    string `json:"path"`
			Context string `json:"context"`
		} `json:"kubernetes"`
		Values map[string]interface{} `json:"values"`
	}
	MerlinConfig struct {
		Environments       []ConfigEnvironment `json:"environments" yaml:"environments"`
		Codefresh          []ConfigCodefresh   `json:"codefresh" yaml:"codefresh"`
		Kubernetes         []ConfigKubernetes  `json:"kubernetes" yaml:"kubernetes"`
		CurrentEnvironment string              `json:"current-environment" yaml:"current-environment"`
	}

	ConfigEnvironment struct {
		Name               string                 `json:"name" yaml:"name"`
		DescriptorLocation string                 `json:"descriptor-location" yaml:"descriptor-location"`
		Codefresh          string                 `json:"codefresh" yaml:"codefresh"`
		Kubernetes         string                 `json:"kubernetes" yaml:"kubernetes"`
		Values             map[string]interface{} `json:"values" yaml:"values"`
	}

	ConfigCodefresh struct {
		Name    string `json:"name" yaml:"name"`
		Path    string `json:"path" yaml:"path"`
		Context string `json:"context" yaml:"context"`
	}

	ConfigKubernetes struct {
		Name    string `json:"name" yaml:"name"`
		Path    string `json:"path" yaml:"path"`
		Context string `json:"context" yaml:"context"`
	}
)

func (m *ActiveConfig) ToJSON() map[string]interface{} {
	res := map[string]interface{}{}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error marshalling: %s\n", err.Error())
		return res
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Printf("Error unmarshalling: %s\n", err.Error())
		return res
	}
	return res
}

func (m *MerlinConfig) BuildActive(name string) (*ActiveConfig, error) {
	if name == "" {
		name = m.CurrentEnvironment
	}
	ac := &ActiveConfig{}
	ac.Name = name

	foundEnv := false
	for _, e := range m.Environments {
		if e.Name == name {
			foundEnv = true
			ac.EnvironmentJS = e.DescriptorLocation
			ac.Values = e.Values
		}
	}
	if !foundEnv {
		return nil, fmt.Errorf("Environment %s not found", name)
	}

	foundCodefresh := false
	for _, c := range m.Codefresh {
		if c.Name == name {
			foundCodefresh = true
			ac.Codefresh.Path = c.Path
			ac.Codefresh.Context = c.Context
		}
	}
	if !foundCodefresh {
		return nil, fmt.Errorf("Codefresh %s not found", name)
	}

	foundKube := false
	for _, k := range m.Kubernetes {
		if k.Name == name {
			foundKube = true
			ac.Kubernetes.Path = k.Path
			ac.Kubernetes.Context = k.Context
		}
	}
	if !foundKube {
		return nil, fmt.Errorf("Kubernetes %s not found", name)
	}

	return ac, nil
}
