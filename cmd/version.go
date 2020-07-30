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

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print merlin version",
	Run: func(cmd *cobra.Command, args []string) {
		v := map[string]interface{}{
			"Version": version,
			"Commit":  commit,
			"Date":    date,
		}
		b, err := json.Marshal(v)
		dieOnError("Failed to marshal struct", err)
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
