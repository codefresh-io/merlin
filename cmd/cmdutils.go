package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/codefresh-io/merlin/pkg/js"
	signalHandler "github.com/codefresh-io/merlin/pkg/signal"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/sirupsen/logrus"
)

func dieIfError(logger *logrus.Entry, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func createSignalHandler(logger *logrus.Entry) signalHandler.Handler {
	return signalHandler.NewSignalhandler([]syscall.Signal{syscall.SIGTERM, syscall.SIGINT}, logger)
}

func resolvePathOrDie(logger *logrus.Entry, path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	usr, err := user.Current()
	dieIfError(logger, err)
	dir := usr.HomeDir
	wd, err := os.Getwd()
	dieIfError(logger, err)

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}

	if strings.HasPrefix(path, ".") {
		abs, err := filepath.Abs(path)
		dieIfError(logger, err)
		return abs
	}

	return filepath.Join(wd, path)
}

func readMerlinEnvironmentFileOrDie(logger *logrus.Entry, path string) *spec.Environment {

	env := &spec.Environment{}
	logger.Debugf("Reading environment.js from %s", path)
	e, err := ioutil.ReadFile(path)
	dieIfError(logger, err)

	jsEngine := js.NewJSEngine()

	_, err = jsEngine.Load([]string{string(e)}, nil)
	dieIfError(logger, err)

	logger.Debugf("Calling build function")
	res, err := jsEngine.CallFn("build", nil)
	dieIfError(logger, err)

	logger.Debugf("Unmarshalling result")
	err = json.Unmarshal([]byte(res.String()), env)
	dieIfError(logger, err)
	return env

}

func getConfig(logger *logrus.Entry, path string, configName string) (*spec.MerlinConfig, error) {
	if path == "" {
		logger.Debug("Path is not passed, using default")
		path = fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
	}
	cnf, err := utils.GetConfigFile(path)
	return cnf, err
}

func getActiveConfig(logger *logrus.Entry, path string, configName string) (*spec.ActiveConfig, error) {
	cnf, err := getConfig(logger, path, configName)
	if err != nil {
		return nil, err
	}
	logger.Debug("Config file found, build active config")
	return cnf.BuildActive(configName)

}
