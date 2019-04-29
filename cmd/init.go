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

	codefresh "github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/codefresh-io/merlin/pkg/kube"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/strvals"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/imdario/mergo"
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
	Run: func(cmd *cobra.Command, args []string) {
		location := fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command":  "Run",
				"Location": location,
			},
			Debug: verbose,
		})
		logger.Debugf("Starting init process - %v", initCmdOpt)
		cnf, err := utils.GetConfigFile(location)
		dieIfError(logger, err)

		cf := spec.ConfigCodefresh{}
		kubernetes := spec.ConfigKubernetes{}
		env := spec.ConfigEnvironment{}

		name := initCmdOpt.name
		if name == "" {
			name = getFromUserOrDie(logger, "Give your environment a name", nil)
		}

		// init references
		{
			cf.Name = name
			kubernetes.Name = name
			env.Name = name
			env.Codefresh = name
			env.Kubernetes = name
		}

		if initCmdOpt.codefresh.path == "" {
			initCmdOpt.codefresh.path = getFromUserOrDie(logger, "Path to codefresh config", nil)
		}
		cf.Path = resolvePathOrDie(logger, initCmdOpt.codefresh.path)

		if initCmdOpt.codefresh.context == "" {
			cf, err := codefresh.GetCFConfig(cf.Path)
			dieIfError(logger, err)
			items := []string{}
			for _, c := range cf.Contexts {
				items = append(items, c.Name)
			}
			name := getFromUserOrDie(logger, "Select codefresh context", items)
			if name == "" {
				name = cf.CurrentContext
			}

			initCmdOpt.codefresh.context = name
		}
		cf.Context = initCmdOpt.codefresh.context

		if initCmdOpt.kubernetes.path == "" {
			initCmdOpt.kubernetes.path = getFromUserOrDie(logger, "Path to kubeconfig", nil)
		}
		kubernetes.Path = resolvePathOrDie(logger, initCmdOpt.kubernetes.path)

		if initCmdOpt.kubernetes.context == "" {
			items, current, err := kube.GetKubeContexts(kubernetes.Path)
			dieIfError(logger, err)
			name := getFromUserOrDie(logger, "Select kube context", items)
			if name == "" {
				name = current
			}
			initCmdOpt.kubernetes.context = name
		}
		kubernetes.Context = initCmdOpt.kubernetes.context

		if initCmdOpt.environmentJs == "" {
			initCmdOpt.environmentJs = getFromUserOrDie(logger, "Path to environment.js file", nil)
		}
		env.DescriptorLocation = resolvePathOrDie(logger, initCmdOpt.environmentJs)

		values := map[string]interface{}{}
		for _, value := range initCmdOpt.setList {
			if err := strvals.ParseInto(value, values); err != nil {
				dieIfError(logger, fmt.Errorf("failed parsing --set data: %s", err))
			}
		}

		err = mergo.Merge(&env.Values, values, mergo.WithOverride, mergo.WithAppendSlice)
		if err != nil {
			dieIfError(logger, err)
		}

		cnf.CurrentEnvironment = name
		cnf.Environments = append(cnf.Environments, env)
		cnf.Codefresh = append(cnf.Codefresh, cf)
		cnf.Kubernetes = append(cnf.Kubernetes, kubernetes)
		err = utils.PersistConfigFile(cnf, fmt.Sprintf("%s/.merlin", os.Getenv("HOME")), "config.yaml")
		dieIfError(logger, err)
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
