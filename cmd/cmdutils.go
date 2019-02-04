package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/codefresh-io/merlin/pkg/cache"
	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ensureLocalDirectory(c *config.Config) error {
	home := os.Getenv("HOME")
	pathToEnsure := fmt.Sprintf("%s/.codefresh/.dev", home)
	if _, err := os.Stat(pathToEnsure); os.IsNotExist(err) {
		if err = os.Mkdir(pathToEnsure, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func dieIfError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func readConfigFromPathOrDie(log *logrus.Entry) *config.Config {
	c := &config.Config{}

	if merlinconfig == "" {
		viper.SetConfigName(config.MerlinFileName)
		viper.AddConfigPath(fmt.Sprintf("$HOME/%s", config.MerlinFilePath))
		viper.SetConfigType(config.MerlinFileExtention)
		err := viper.ReadInConfig()
		dieIfError(err)
		err = viper.Unmarshal(c)
		dieIfError(err)
	} else {
		log.Debugf("Reading file %s", merlinconfig)
		err := utils.ReadFileInto(merlinconfig, c)
		dieIfError(err)
	}
	return c
}

func createConfigFile(c *config.Config, merlinconfig string) error {
	var err error
	home := os.Getenv("HOME")
	if c.Environment.Spec.Path != "" {
		c.Environment.Spec.Path, err = filepath.Abs(c.Environment.Spec.Path)
	} else {
		g := c.Environment.Spec.Git
		g.Owner = "codefresh-io"
		g.Repo = "cf-helm"
		g.Path = "codefresh/env/dynamic/environment.yaml.tmpl"
		g.Revision = "master"
	}
	dieIfError(err)

	if c.Codefresh.Path == "" {
		c.Codefresh.Path, err = filepath.Abs(fmt.Sprintf("%s/.cfconfig", home))
	} else {
		c.Codefresh.Path, err = filepath.Abs(c.Codefresh.Path)
	}
	dieIfError(err)

	if c.Kube.Path == "" {
		c.Kube.Path, err = filepath.Abs(fmt.Sprintf("%s/.kube/config", home))
	} else {
		c.Kube.Path, err = filepath.Abs(c.Kube.Path)
	}
	dieIfError(err)

	var filePath string
	if merlinconfig == "" {
		filePath = fmt.Sprintf("%s/%s/%s.%s", home, config.MerlinFilePath, config.MerlinFileName, config.MerlinFileExtention)
	} else {
		filePath, err = filepath.Abs(merlinconfig)
		dieIfError(err)
	}
	res, err := yaml.Marshal(c)
	dieIfError(err)
	err = ioutil.WriteFile(filePath, res, 0644)
	return err
}

func createCacheStore(c *config.Config, noCache bool, log *logrus.Entry) cache.Cache {
	return cache.New(&cache.Options{
		Path:    c.Cache.Path,
		Logger:  log,
		NoCache: noCache,
	})
}
