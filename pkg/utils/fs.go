package utils

import (
	"fmt"
	"github.com/codefresh-io/merlin/pkg/spec"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetConfigFile(path string) (*spec.MerlinConfig, error) {
	cnf := &spec.MerlinConfig{}
	f, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		dir, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
		return cnf, nil
	}
	err = yaml.Unmarshal(f, cnf)
	return cnf, err
}

func PersistConfigFile(config *spec.MerlinConfig, path string, name string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, name), b, os.ModePerm); err != nil {
		return err
	}
	return nil
}
