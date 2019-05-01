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
	"github.com/spf13/cobra"
)

// configList represents the version command
var configList = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "config-list",
			},
			Debug: verbose,
		})
		cnf, err := getConfig(logger, "", "")
		dieIfError(logger, err)
		for _, e := range cnf.Environments {
			name := e.Name
			ac, err := cnf.BuildActive(e.Name)
			dieIfError(logger, err)
			if e.Name == cnf.CurrentEnvironment {
				name = fmt.Sprintf("%s (active)", name)
			}
			json := ac.ToJSONString()
			fmt.Printf("Name: %s\n", name)
			fmt.Printf("JSON: %s\n", json)
			fmt.Println()
		}
	},
}

func init() {
	configCmd.AddCommand(configList)
}
