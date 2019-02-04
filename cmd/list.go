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
	"github.com/codefresh-io/merlin/pkg/environment"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/codefresh-io/merlin/pkg/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show a list of all commands exposed",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "List",
			},
			Debug: verbose,
		})
		c := readConfigFromPathOrDie(log)
		store := createCacheStore(c, noCache, log)
		defer store.Persist()
		res, err := environment.Build(c, store, log).List(nil)
		dieIfError(err)
		t := table.New(&table.Options{
			Headers: []string{"Type", "Level", "Name", "Description"},
		})
		t.Table().AppendBulk(res)
		t.Table().Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}