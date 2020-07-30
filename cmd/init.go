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
	"fmt"
	"os"
	"os/user"

	"github.com/codefresh-io/go/logger"
	"github.com/codefresh-io/merlin/pkg/kube"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmdOpt struct {
	name       string
	kubernetes struct {
		path    string
		context string
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	Run: func(cmd *cobra.Command, args []string) {
		location := fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
		logger, err := logger.New(logger.Options{
			Context: []interface{}{
				"Command", "Run",
				"Location", location,
			},
		})
		dieOnError("Failed to create logger", err)
		logger.Debug("Starting init process", "cmd", initCmdOpt)

		name := initCmdOpt.name
		if name == "" {
			defaultEnvName := ""
			u, _ := user.Current()
			if u != nil {
				defaultEnvName = u.Username
			}
			name = getFromUserOrDie(logger, "Set environment name", nil, defaultEnvName)
		}

		if initCmdOpt.kubernetes.path == "" {
			initCmdOpt.kubernetes.path = getFromUserOrDie(logger, "Set path to Kubernetes config file (kubeconfig)", nil, fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")))
		}
		kubePath := resolvePathOrDie(logger, initCmdOpt.kubernetes.path)

		if initCmdOpt.kubernetes.context == "" {
			items, current, err := kube.GetKubeContexts(kubePath)
			dieOnError("Failed to get context from kubeconfig file", err)
			name := getFromUserOrDie(logger, "Select Kubernetes context to be used", items, "")
			if name == "" {
				name = current
			}
			initCmdOpt.kubernetes.context = name
		}
		kubeContext := initCmdOpt.kubernetes.context
		kubenamespace := getFromUserOrDie(logger, "Namespace", nil, name)
		shell := getFromUserOrDie(logger, "What is your shell?", nil, "bash")
		cnf := &spec.Config{
			Name: name,
			Cluster: spec.Cluster{
				Context:   kubeContext,
				Namespace: kubenamespace,
				Path:      kubePath,
			},
			Shell: shell,
		}
		err = utils.PersistConfigFile(cnf, fmt.Sprintf("%s/.merlin", os.Getenv("HOME")), "config.yaml")
		dieOnError("Failed to persist config file", err)
	},
}

func getFromUserOrDie(logger logger.Logger, label string, items []string, defaultValue string) string {
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
			Label:   label,
			Default: defaultValue,
		}
		res, err = p.Run()
	}
	dieOnError("Failed to get input", err)
	return res
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCmdOpt.name, "name", viper.GetString("MERLIN_ENVIRONMENT_NAME"), "Set environment name [$MERLIN_ENVIRONMENT_NAME]")
	initCmd.Flags().StringVar(&initCmdOpt.kubernetes.path, "kubernetes-config-path", viper.GetString("KUBECONFIG"), "")
	initCmd.Flags().StringVar(&initCmdOpt.kubernetes.context, "kubernetes-config-context", viper.GetString("KUBE_CONTEXT"), "")
}
