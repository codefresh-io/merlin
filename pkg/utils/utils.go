package utils

import (
	"io/ioutil"
	"net"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

func GetAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}

func ReadFileInto(path string, target interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, target)
	return err
}
