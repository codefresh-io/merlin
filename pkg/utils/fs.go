package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codefresh-io/merlin/pkg/spec"
	"gopkg.in/yaml.v2"
)

func GetConfigFile(path string) (*spec.Config, error) {
	cnf := &spec.Config{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, cnf)
	return cnf, err
}

func PersistConfigFile(config *spec.Config, path string, name string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, name), b, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func GetSerivceFile(path string) (*spec.Service, error) {
	cnf := &spec.Service{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, cnf)
	return cnf, err
}
