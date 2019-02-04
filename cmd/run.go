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
	"errors"

	"github.com/codefresh-io/merlin/pkg/environment"
	"github.com/codefresh-io/merlin/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	set             []string
	component       string
	skipCommandExec bool
)

var runCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run command",
	Long:  "merlin run connect --component cfapi",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			dieIfError(errors.New("Run command from component template"))
		}
		if len(args) > 1 {
			dieIfError(errors.New("Cant run multiple commands"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "Run",
				"Target":  args[0],
			},
			Debug: verbose,
		})
		c := readConfigFromPathOrDie(log)
		store := createCacheStore(c, noCache, log)
		defer store.Persist()
		err := environment.Build(c, store, log).Run(&environment.RunCommandOptions{
			Component: component,
			Override:  set,
			Operator:  args[0],
			SkipExec:  skipCommandExec,
		})
		dieIfError(err)
	},
}

func init() {
	viper.BindEnv("component", "MERLIN_COMPONENT")

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringArrayVar(&set, "set", nil, "Set value to override")
	runCmd.Flags().StringVar(&component, "component", viper.GetString("component"), "Set name of the component where the command exist [$MERLIN_COMPONENT]")
	runCmd.Flags().BoolVar(&skipCommandExec, "skip-exec", false, "set to get the final command as string, without executing it")
}
