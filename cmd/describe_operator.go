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
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var describeOperatorCmdOpt struct {
	environment  string
	merlinconfig string
}

var describeOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Show a list of all operators exposed",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "List",
			},
			Debug: verbose,
		})
		if len(args) > 1 {
			dieIfError(logger, fmt.Errorf("Passed too many operators"))
		}
		if len(args) == 0 {
			dieIfError(logger, fmt.Errorf("no operator passed"))
		}
		ac, err := getConfig(logger, describeEnvCmdOpt.merlinconfig, describeEnvCmdOpt.environment)
		dieIfError(logger, err)

		env := readMerlinEnvironmentFileOrDie(logger, ac.EnvironmentJS)

		logger.Debug("Printing table")
		t := table.New(&table.Options{
			Headers: []string{"Name", "Description", "Environment Variables", "Default", "Required", "Allow Interactove"},
		})
		op := spec.Operator{}
		for _, o := range env.Operators {
			if o.Name == args[0] {
				op = o
			}
		}
		if &op == nil {
			dieIfError(logger, fmt.Errorf("Operator %s not found", args[0]))
		}
		scope := op.Scope
		if scope == "" {
			scope = "environment"
		}
		fmt.Printf("Operator: %s (%s)\n", op.Name, scope)
		fmt.Println(op.Description)

		if len(op.Params) > 0 {
			for _, p := range op.Params {
				description := p.Description
				if len(p.Description) > 30 {
					description = fmt.Sprintf("%s...", p.Description[:30])
				}
				t.Table().Append([]string{
					p.Name,
					description,
					p.EnvironmentVariable,
					p.Default,
					fmt.Sprintf("%t", p.Required),
					fmt.Sprintf("%t", p.InteractiveInput),
				})
			}
			t.Table().Render()
		}
	},
}

func init() {
	describeCmd.AddCommand(describeOperatorCmd)
	describeOperatorCmd.Flags().StringVar(&describeOperatorCmdOpt.merlinconfig, "merlinconfig", viper.GetString("MERLIN_CONFIG"), "Path to merlinconfig file (default $HOME/.merlin/config) [$MERLIN_CONFIG]")
	describeOperatorCmd.Flags().StringVar(&describeOperatorCmdOpt.environment, "environment", viper.GetString("MERLIN_ENVIRONMENT"), "Name of the environment from merlinconfig [$MERLIN_ENVIRONMENT]")
}
