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

	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/spf13/cobra"
)

// configUseCmd represents the version command
var configUseCmd = &cobra.Command{
	Use: "use",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "config-use",
			},
			Debug: verbose,
		})
		if len(args) > 1 {
			dieIfError(logger, fmt.Errorf("Passed too many arguments"))
		}
		if len(args) == 0 {
			dieIfError(logger, fmt.Errorf("no argument passed"))
		}

		cnf, err := getConfig(logger, "", args[0])
		dieIfError(logger, err)
		_, err = cnf.BuildActive(args[0])
		dieIfError(logger, err)
		logger.Debugf("Environment %s found, setting to be the current one", args[0])
		cnf.CurrentEnvironment = args[0]
		err = utils.PersistConfigFile(cnf, fmt.Sprintf("%s/.merlin", os.Getenv("HOME")), "config.yaml")
		dieIfError(logger, err)
		logger.Infof("Current Environment updated to %s", args[0])
	},
}

func init() {
	configCmd.AddCommand(configUseCmd)
}
