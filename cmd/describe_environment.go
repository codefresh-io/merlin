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

	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var describeEnvCmdOpt struct {
	environment string
}

var describeEnvCmd = &cobra.Command{
	Use:   "environment",
	Short: "Show a list of all operators exposed",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "List",
			},
			Debug: verbose,
		})
		env := readMerlinEnvironmentFileOrDie(logger, describeEnvCmdOpt.environment)

		logger.Debug("Printing table")
		t := table.New(&table.Options{
			Headers: []string{"Operator", "Description"},
		})
		for _, o := range env.Operators {
			if o.Scope == "" {
				o.Scope = "environment"
			}
			description := o.Description
			if len(o.Description) > 30 {
				description = fmt.Sprintf("%s...", o.Description[:30])
			}
			t.Table().Append([]string{
				fmt.Sprintf("%s (%s)", o.Name, o.Scope),
				description,
			})
		}
		t.Table().Render()
		return nil
	},
}

func init() {
	describeCmd.AddCommand(describeEnvCmd)
	describeEnvCmd.Flags().StringVar(&describeEnvCmdOpt.environment, "environment", viper.GetString("MERLIN_ENVIRONMENT"), "Paht to environment.js file [$MERLIN_ENVIRONMENT]")
	describeEnvCmd.MarkFlagRequired("environment")
}
