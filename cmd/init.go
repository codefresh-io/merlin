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
	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/github"
	"github.com/codefresh-io/merlin/pkg/logger"

	"github.com/spf13/cobra"
)

var c = &config.Config{}

var configCmd = &cobra.Command{
	Use:   "init",
	Short: "Create config file",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "Create",
			},
			Debug: verbose,
		})
		c.Name = c.Kube.Namespace
		if merlinconfig == "" {
			err := ensureLocalDirectory(c)
			dieIfError(err)
		}
		err := createConfigFile(c, merlinconfig)
		dieIfError(err)
		if c.Environment.Path != "" {
			_, err = config.ReadEnvironmentDescriptor(c.Environment.Path)
			dieIfError(err)
		} else {
			g := c.Environment.Git
			git, err := github.New(c.Github.Token, log)
			dieIfError(err)
			_, err = git.ReadFile(g.Owner, g.Repo, g.Path, g.Revision)
			dieIfError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	//make sure the path has relative prefix, or add it (--path codefresh.yaml will not have full path)
	configCmd.Flags().StringVar(&c.Codefresh.Context, "codefresh-config-context", "", "Set name the context name in codefresh config to use")
	configCmd.Flags().StringVar(&c.Codefresh.Path, "codefresh-config-path", "", "Set path to codefresh config")
	configCmd.Flags().StringVar(&c.Kube.Context, "kube-config-context", "", "Set name the in kubeconfig to use")
	configCmd.Flags().StringVar(&c.Kube.Path, "kube-config-path", "", "Set path to kubeconfig")
	configCmd.Flags().StringVar(&c.Kube.Namespace, "kube-config-namespace", "", "Set name for the environment")
	configCmd.Flags().StringVar(&c.Github.Token, "github-token", "", "Set token to github")
	configCmd.Flags().StringVar(&c.Environment.Path, "environment-descriptor", "", "Set path to environment descriptor")
}
