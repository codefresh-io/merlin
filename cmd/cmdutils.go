package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/codefresh-io/go/logger"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/sirupsen/logrus"
)

func dieOnError(message string, err error) {
	if err != nil {
		fmt.Printf("[ERROR]: %s - %v", message, err)
		os.Exit(1)
	}
}

func resolvePathOrDie(logger logger.Logger, path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	usr, err := user.Current()
	dieOnError("Failed to resolve current user", err)
	dir := usr.HomeDir
	wd, err := os.Getwd()
	dieOnError("Failed to resolve working directory", err)

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}

	if strings.HasPrefix(path, ".") {
		abs, err := filepath.Abs(path)
		dieOnError("Failed to get absolute filepath", err)
		return abs
	}

	return filepath.Join(wd, path)
}

func getConfig(logger *logrus.Entry, path string) (*spec.Config, error) {
	if path == "" {
		logger.Debug("Path is not passed, using default")
		path = fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
	}
	cnf, err := utils.GetConfigFile(path)
	return cnf, err
}

func startSignalHandler(logger logger.Logger, cancel context.CancelFunc) {
	var sig = make(chan os.Signal)
	go func() {
		logger.Debug("Waiting to get signal")
		<-sig
		cancel()
	}()

}
