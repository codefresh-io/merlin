package cmd

/*
Copyright 2019 The Codefresh Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	codefresh "github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/codefresh-io/merlin/pkg/kube"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmdOpt struct {
	name          string
	setList       []string
	environmentJs string
	codefresh     struct {
		path    string
		context string
	}
	kubernetes struct {
		path    string
		context string
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "Run",
			},
			Debug: verbose,
		})
		logger.Debug("Running")

		cnf := spec.MerlinConfig{}

		if initCmdOpt.name == "" {
			cnf.Name = getFromUserOrDie(logger, "Give your environment a name", nil)
		} else {
			cnf.Name = initCmdOpt.name
		}
		location := fmt.Sprintf("%s/.merlin/%s.json", os.Getenv("HOME"), cnf.Name)

		if initCmdOpt.codefresh.path == "" {
			path := getFromUserOrDie(logger, "Path to codefresh config", nil)
			cnf.Codefresh.Path = resolvePathOrDie(logger, path)
		} else {
			cnf.Codefresh.Path = initCmdOpt.codefresh.path
		}
		if initCmdOpt.codefresh.context == "" {
			cf, err := codefresh.GetCFConfig(cnf.Codefresh.Path)
			dieIfError(logger, err)
			items := []string{}
			for _, c := range cf.Contexts {
				items = append(items, c.Name)
			}
			name := getFromUserOrDie(logger, "Select codefresh context", items)
			if name == "" {
				name = cf.CurrentContext
			}

			cnf.Codefresh.Context = name
		} else {
			cnf.Codefresh.Context = initCmdOpt.codefresh.context
		}

		if initCmdOpt.kubernetes.path == "" {
			path := getFromUserOrDie(logger, "Path to kubeconfig", nil)
			cnf.Kubernetes.Path = resolvePathOrDie(logger, path)
		} else {
			cnf.Kubernetes.Path = initCmdOpt.kubernetes.path
		}

		if initCmdOpt.kubernetes.context == "" {
			items, current, err := kube.GetKubeContexts(cnf.Kubernetes.Path)
			dieIfError(logger, err)
			name := getFromUserOrDie(logger, "Select kube context", items)
			if name == "" {
				name = current
			}

			cnf.Kubernetes.Context = name
		} else {
			cnf.Kubernetes.Context = initCmdOpt.kubernetes.context
		}

		if initCmdOpt.environmentJs == "" {
			path := getFromUserOrDie(logger, "Path to environment.js file", nil)
			cnf.EnvironmentJS = resolvePathOrDie(logger, path)
		} else {
			cnf.EnvironmentJS = initCmdOpt.environmentJs
		}

		res, err := json.Marshal(cnf)
		dieIfError(logger, err)
		err = ioutil.WriteFile(location, res, 0644)
		return err

	},
}

func getFromUserOrDie(logger *logrus.Entry, label string, items []string) string {
	var res string
	var err error
	if items != nil && len(items) > 0 {
		p := promptui.Select{
			Items: items,
			Label: label,
		}
		_, res, err = p.Run()
	} else {
		p := promptui.Prompt{
			Label: label,
		}
		res, err = p.Run()
	}
	dieIfError(logger, err)
	return res
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringArrayVar(&initCmdOpt.setList, "set", []string{}, "--set name=value OR --set key.inner_key=value")

	initCmd.Flags().StringVar(&initCmdOpt.name, "name", viper.GetString("MERLIN_ENVIRONMENT_NAME"), "Set environment name [$MERLIN_ENVIRONMENT_NAME]")
	initCmd.Flags().StringVar(&initCmdOpt.environmentJs, "environment", viper.GetString("MERLIN_ENVIRONMENT_JS_PATH"), "Set path to environment.js file [$MERLIN_ENVIRONMENT_JS_PATH]")

	initCmd.Flags().StringVar(&initCmdOpt.codefresh.path, "codefresh-config-path", viper.GetString("CODEFRESH_CONFIG"), "")
	initCmd.Flags().StringVar(&initCmdOpt.codefresh.context, "codefresh-config-context", viper.GetString("CODEFRESH_CONTEXT"), "")

	initCmd.Flags().StringVar(&initCmdOpt.kubernetes.path, "kubernetes-config-path", viper.GetString("KUBECONFIG"), "")
	initCmd.Flags().StringVar(&initCmdOpt.kubernetes.context, "kubernetes-config-context", viper.GetString("KUBE_CONTEXT"), "")
}
