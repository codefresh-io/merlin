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
	"context"
	"fmt"
	"os"
	"path"

	"github.com/codefresh-io/go/logger"
	"github.com/codefresh-io/merlin/pkg/commander"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.New(logger.Options{
			Context: []interface{}{"cmd", "connect"},
		})
		dieOnError("Failed to create logger", err)
		ctx, cancel := context.WithCancel(context.Background())
		go startSignalHandler(logger, cancel)

		location := fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
		cnf, err := getConfig(nil, location)
		dieOnError("Failed to read environment config file", err)
		logger.Debug("Connecting", "env", cnf.Name)
		pwd, err := os.Getwd()
		dieOnError("Failed to get current working dir", err)
		logger.Debug("Loading service.yaml file")
		svc, err := utils.GetSerivceFile(path.Join(pwd, "service.yaml"))
		dieOnError("Failed to read service.yaml file", err)
		logger.Debug("Serivce loaded")
		tpArgs := []string{
			"--swap-deployment",
			fmt.Sprintf("%s-%s-base", cnf.Name, svc.Name),
			"--context", cnf.Cluster.Context,
			"--namespace", cnf.Cluster.Namespace,
		}
		if cnf.Shell != "" {
			tpArgs = append(tpArgs, []string{"--run", cnf.Shell}...)
		}
		tpEnv := []string{}
		for _, p := range svc.Ports {
			port, err := utils.GetAvailablePort()
			dieOnError("Failed to generate port", err)
			tpArgs = append(tpArgs, []string{"--expose", fmt.Sprintf("%d:%d", port, p.Default)}...)
			tpEnv = append(tpEnv, fmt.Sprintf("merlin_generated_%s=%d", p.EnvVar, port))
		}

		for _, e := range svc.Environment {
			tpEnv = append(tpEnv, fmt.Sprintf("merlin_generated_%s=%s", e.Name, e.Default))
		}

		opt := &commander.Options{
			Program:  "telepresence",
			Detached: true,
			Logger:   logger,
			Args:     tpArgs,
			WorkDir:  pwd,
			Env:      tpEnv,
		}

		tpCmd := commander.New(opt)
		fmt.Println(tpCmd.DryRun())
		if err := tpCmd.Run(ctx); err != nil {
			dieOnError("Failed to run connect command", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
