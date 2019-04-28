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

	"github.com/codefresh-io/merlin/pkg/commander"
	"github.com/codefresh-io/merlin/pkg/js"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/strvals"
	"github.com/imdario/mergo"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmdOpt struct {
	merlinconfigPath string
	dryRun           bool
	setList          []string
	componentName    string
	codefresh        struct {
		path    string
		context string
	}
	kubernetes struct {
		path      string
		context   string
		namespace string
	}
	operator    *spec.Operator
	component   *spec.Component
	environemnt *spec.Environment
	config      *spec.MerlinConfig
}

var testCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a operator",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "Run",
			},
			Debug: verbose,
		})
		logger.Debug("Running")
		if runCmdOpt.merlinconfigPath != "" {
			logger.Debugf("Reading merlinconfig file from %s", runCmdOpt.merlinconfigPath)
			configFile, err := ioutil.ReadFile(resolvePathOrDie(logger, runCmdOpt.merlinconfigPath))
			if !os.IsNotExist(err) {
				runCmdOpt.config = &spec.MerlinConfig{}
				err = json.Unmarshal(configFile, runCmdOpt.config)
				if err != nil {
					dieIfError(logger, err)
				}
			}
		}

		logger.Debugf("Reading environment.js file from %s", runCmdOpt.config.EnvironmentJS)
		f, err := ioutil.ReadFile(runCmdOpt.config.EnvironmentJS)
		if err != nil {
			dieIfError(logger, err)
		}
		logger.Debugf("Creating new JS engine")
		jsEngine := js.NewJSEngine()

		_, err = jsEngine.Load([]string{string(f)}, nil)
		if err != nil {
			dieIfError(logger, err)
		}
		logger.Debugf("Calling build function")
		res, err := jsEngine.CallFn("build", nil)
		if err != nil {
			dieIfError(logger, err)
		}
		logger.Debugf("Unmarshalling result")
		runCmdOpt.environemnt = &spec.Environment{}
		err = json.Unmarshal([]byte(res.String()), runCmdOpt.environemnt)
		if err != nil {
			dieIfError(logger, err)
		}

		logger.Debugf("Looking for operator %s in environment.js", args[0])
		for _, o := range runCmdOpt.environemnt.Operators {
			if o.Name == args[0] {
				op := o
				runCmdOpt.operator = &op
			}
		}

		if runCmdOpt.operator == nil {
			dieIfError(logger, fmt.Errorf("Operator %s not found", args[0]))
		}

		params := map[string]interface{}{}
		for _, p := range runCmdOpt.operator.Params {
			var res = ""
			logger.Debugf("Preparing parameter %s, description: %s", p.Name, p.Description)
			if p.Required {
				logger.Debugf("Parameter is required")
			}

			if p.Default != "" {
				logger.Debugf("Default value foundf")
				res = p.Default
			}

			if p.EnvironmentVariable != "" {
				logger.Debugf("Checking environment variable to have %s", p.EnvironmentVariable)
				r := os.Getenv(p.EnvironmentVariable)
				if r != "" {
					logger.Debugf("Found param in environment variables")
					res = r
				}
			}

			if p.InteractiveInput && res == "" {
				label := ""
				if p.Description != "" {
					label = p.Description
				} else {
					label = p.Name
				}
				prompt := promptui.Prompt{
					Label: label,
				}
				result, err := prompt.Run()
				if err != nil {
					dieIfError(logger, err)
				}
				res = result
			}
			if res == "" && p.Required {
				dieIfError(logger, fmt.Errorf("Error: Parameter %s is required by operator %s and not been found", p.Name, runCmdOpt.operator.Name))
			}
			params[p.Name] = res
		}

		if runCmdOpt.componentName != "" {
			logger.Debugf("Searching for component %s in environemnt.js", runCmdOpt.componentName)
			for _, c := range runCmdOpt.environemnt.Components {
				logger.Debugf("Matching component %s=%s", c.Name, runCmdOpt.componentName)
				if c.Name == runCmdOpt.componentName {
					logger.Debugf("Found component %s", c.Name)
					com := c
					runCmdOpt.component = &com
				}
			}
			if runCmdOpt.component == nil {
				dieIfError(logger, fmt.Errorf("Component %s was passed but not found", args[0]))
			}
		}
		if runCmdOpt.operator.Scope == "component" && runCmdOpt.component == nil {
			dieIfError(logger, fmt.Errorf("Operator \"%s\" is a scoped component, but no component was spesified, use --component", runCmdOpt.operator.Name))
		}

		logger.Debugf("Executing operator %s", runCmdOpt.operator.Name)
		input := runCmdOpt.config.ToJSON()
		err = mergo.Merge(&input, params, mergo.WithOverride, mergo.WithAppendSlice)
		if err != nil {
			dieIfError(logger, err)
		}

		override := map[string]interface{}{}
		for _, value := range runCmdOpt.setList {
			if err := strvals.ParseInto(value, override); err != nil {
				dieIfError(logger, fmt.Errorf("failed parsing --set data: %s", err))
			}
		}
		err = mergo.Merge(&input, override, mergo.WithOverride, mergo.WithAppendSlice)
		if err != nil {
			dieIfError(logger, err)
		}

		execArr, err := jsEngine.CallFn(fmt.Sprintf("$%s", runCmdOpt.operator.Name), input, runCmdOpt.component.ToJSON())
		if err != nil {
			dieIfError(logger, err)
		}

		set := spec.CmdSet{}
		err = json.Unmarshal([]byte(execArr.String()), &set)
		if err != nil {
			dieIfError(logger, err)
		}

		for _, c := range set {
			logger.Debugf("cmd: %v", c)
			if !runCmdOpt.dryRun {
				cmd := commander.New(&commander.Options{
					Program:  c.Exec[0],
					Args:     c.Exec[1:],
					Env:      c.Env,
					Detached: c.Detached,
					Logger:   logger,
					WorkDir:  resolvePathOrDie(logger, c.WorkDir),
				})
				err = cmd.Run()
				if err != nil {
					dieIfError(logger, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVar(&runCmdOpt.componentName, "component", viper.GetString("MERLIN_COMPONENT"), "Name of the component to be execute as part of the operator [$MERLIN_COMPONENT]")
	testCmd.Flags().StringVar(&runCmdOpt.merlinconfigPath, "merlinconfig", viper.GetString("MERLIN_CONFIG"), "Path to relevant merlinconfig file (mostly in $HOME/.merlin/) [$MERLIN_CONFIG]")
	testCmd.Flags().StringArrayVar(&runCmdOpt.setList, "set", []string{}, "--set name=value OR --set key.inner_key=value")
	testCmd.Flags().BoolVar(&runCmdOpt.dryRun, "dry-run", false, "Dry run")

	testCmd.Flags().StringVar(&runCmdOpt.codefresh.path, "codefresh-config-path", viper.GetString("CODEFRESH_CONFIG"), "")
	testCmd.Flags().StringVar(&runCmdOpt.codefresh.context, "codefresh-config-context", viper.GetString("CODEFRESH_CONTEXT"), "")

	testCmd.Flags().StringVar(&runCmdOpt.kubernetes.path, "kubernetes-config-path", viper.GetString("KUBECONFIG"), "")
	testCmd.Flags().StringVar(&runCmdOpt.kubernetes.context, "kubernetes-config-context", viper.GetString("KUBE_CONTEXT"), "")
	testCmd.Flags().StringVar(&runCmdOpt.kubernetes.namespace, "kubernetes-config-namespace", viper.GetString("KUBE_NAMESPACE"), "")
}
