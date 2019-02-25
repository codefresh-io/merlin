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

	"fmt"
	"os"

	codefresh "github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/github"
	"github.com/codefresh-io/merlin/pkg/kube"
	"github.com/codefresh-io/merlin/pkg/logger"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var c = &config.Config{}

var configCmd = &cobra.Command{
	Use:     "init [name]",
	Short:   "Create config file",
	PreRunE: runInteractiveShell,
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
		if c.Environment.Spec.Path != "" {
			_, err = config.ReadEnvironmentDescriptor(c.Environment.Spec.Path)
			dieIfError(err)
		} else {
			g := c.Environment.Spec.Git
			git := github.New(c.Github.Token, log)
			dieIfError(err)
			_, err = git.ReadFile(g.Owner, g.Repo, g.Path, g.Revision)
			dieIfError(err)
		}
	},
}

func runInteractiveShell(cmd *cobra.Command, args []string) error {
	if c.Github.Token == "" {
		prompt := promptui.Prompt{
			Label: "Set token to github",
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Github.Token = result
	}

	if c.Codefresh.Path == "" {
		prompt := promptui.Prompt{
			Label:   "Set full path to codefresh config",
			Default: fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME")),
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Codefresh.Path = result
	}

	if c.Codefresh.Context == "" {
		cfconfig, err := codefresh.GetCFConfig(c.Codefresh.Path)
		if err != nil {
			return err
		}
		items := []string{}
		for _, c := range cfconfig.Contexts {
			items = append(items, c.Name)
		}
		prompt := promptui.Select{
			Label: "Set the name of to codefresh context",
			Items: items,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		c.Codefresh.Context = result
	}

	if c.Environment.Spec.Git.Owner == "" {
		prompt := promptui.Prompt{
			Label: "Set the repository owner where environment descriptor can be found",
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Environment.Spec.Git.Owner = result
	}
	if c.Environment.Spec.Git.Repo == "" {
		prompt := promptui.Prompt{
			Label: "Set the repository name where environment descriptor can be found",
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Environment.Spec.Git.Repo = result
	}
	if c.Environment.Spec.Git.Path == "" {
		prompt := promptui.Prompt{
			Label: "Set the path where the environment descriptor can be found relative to the repository",
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Environment.Spec.Git.Path = result
	}

	if c.Kube.Path == "" {
		prompt := promptui.Prompt{
			Label:   "Set path to kubectl config",
			Default: fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
		}
		result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Kube.Path = result
	}
	if c.Kube.Context == "" {
		items, err := kube.GetKubeContexts(c.Kube.Path)
		if err != nil {
			return err
		}
		prompt := promptui.Select{
			Label: "Set the name of the context in kubeconfig",
			Items: items,
		}
		_, result, err := prompt.Run()

		if err != nil {
			return err
		}
		c.Kube.Context = result
	}
	return nil
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

	configCmd.Flags().StringVar(&c.Environment.Spec.Path, "environment-descriptor", "", "Set path to environment descriptor from local machine")
	configCmd.Flags().StringVar(&c.Environment.Spec.Git.Owner, "environment-descriptor-repo-owner", "", "Set the repo owner of your environment descriptor")
	configCmd.Flags().StringVar(&c.Environment.Spec.Git.Repo, "environment-descriptor-repo-name", "", "Set the repo name of your environment descriptor")
	configCmd.Flags().StringVar(&c.Environment.Spec.Git.Path, "environment-descriptor-path", "", "Set the path to the your environment descriptor relative to the repositoty")
	configCmd.Flags().StringVar(&c.Environment.Spec.Git.Revision, "environment-descriptor-revision", "master", "Set the revision of your environment descriptor relative to the repositoty (default is master)")

	rootCmd.MarkFlagRequired("codefresh-config-context")
	rootCmd.MarkFlagRequired("codefresh-config-path")
	rootCmd.MarkFlagRequired("kube-config-context")
	rootCmd.MarkFlagRequired("kube-config-path")
	rootCmd.MarkFlagRequired("github-token")
}
