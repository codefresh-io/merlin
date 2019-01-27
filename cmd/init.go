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

	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/github"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var c = &config.Config{}

var configCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Create config file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			dieIfError(errors.New("Name is required"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c.Name = args[0]
		c.Kube.Namespace = args[0]
		log := logger.New(&logger.LoggerOptions{
			Fields: map[string]interface{}{
				"Command": "Create",
			},
			Debug: verbose,
		})
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

	viper.BindEnv("kubeconfig", "KUBECONFIG")
	viper.BindEnv("cfconfig", "CFCONFIG")
	viper.BindEnv("github", "GITHUB_TOKEN")

	configCmd.Flags().StringVar(&c.Codefresh.Context, "codefresh-config-context", "", "Set name the context name in codefresh config to use (required)")
	configCmd.Flags().StringVar(&c.Codefresh.Path, "codefresh-config-path", viper.GetString("cfconfig"), "Set path to codefresh config [$CFCONFIG]")
	configCmd.Flags().StringVar(&c.Kube.Context, "kube-config-context", "", "Set name the in kubeconfig to use (required)")
	configCmd.Flags().StringVar(&c.Kube.Path, "kube-config-path", viper.GetString("kubeconfig"), "Set path to kubeconfig [$KUBECONFIG] (default: $HOME/.kube/config)")
	configCmd.Flags().StringVar(&c.Github.Token, "github-token", viper.GetString("github"), "Set token to github [$GITHUB_TOKEN]")
	configCmd.Flags().StringVar(&c.Environment.Path, "environment-descriptor", "", "Set path to environment descriptor (default is from github)")

	rootCmd.MarkFlagRequired("codefresh-config-context")
	rootCmd.MarkFlagRequired("codefresh-config-path")
	rootCmd.MarkFlagRequired("kube-config-context")
	rootCmd.MarkFlagRequired("kube-config-path")
	rootCmd.MarkFlagRequired("github-token")
}
