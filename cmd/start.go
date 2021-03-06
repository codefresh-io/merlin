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
	"strconv"
	"strings"

	"github.com/codefresh-io/go/logger"
	"github.com/codefresh-io/merlin/pkg/commander"
	"github.com/codefresh-io/merlin/pkg/spec"
	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.New(logger.Options{
			Context: []interface{}{"cmd", "start"},
		})
		dieOnError("Failed to create logger", err)

		ctx, cancel := context.WithCancel(context.Background())
		go startSignalHandler(logger, cancel)

		logger.Debug("Loading service.yaml file")
		pwd, err := os.Getwd()
		dieOnError("Failed to get current working dir", err)
		svc, err := utils.GetSerivceFile(path.Join(pwd, "service.yaml"))
		dieOnError("Failed to read service.yaml file", err)

		location := fmt.Sprintf("%s/.merlin/config.yaml", os.Getenv("HOME"))
		cnf, err := getConfig(nil, location)
		dieOnError("Failed to read environment config file", err)

		debug := svc.Debug.Port.Default
		if os.Getenv(svc.Debug.Port.EnvVar) != "" {
			p, err := strconv.Atoi(os.Getenv(svc.Debug.Port.EnvVar))
			dieOnError("Failed to conver port to number", err)
			debug = p
		}
		if cnf.PortGenerationStrategy.Debug != nil && *cnf.PortGenerationStrategy.Debug == spec.PortGenerationStrategyRandom {
			d, err := utils.GetAvailablePort()
			dieOnError("Failed to generate port for debug", err)
			debug = d
		}

		tpEnv := []string{
			fmt.Sprintf("DEBUG_PORT=%d", debug),
		}
		for _, env := range os.Environ() {
			kv := strings.Split(env, "=")
			if !strings.HasPrefix(kv[0], "merlin_generated_") {
				continue
			}
			original := strings.Split(kv[0], "merlin_generated_")
			tpEnv = append(tpEnv, fmt.Sprintf("%s=%s", original[1], kv[1]))
		}

		opt := &commander.Options{
			Program:  args[0],
			Detached: true,
			Logger:   logger,
			Args:     args[1:],
			WorkDir:  pwd,
			Env:      tpEnv,
		}
		logger.Debug("Starting service", "name", svc.Name)
		tpCmd := commander.New(opt)
		fmt.Println(tpCmd.DryRun())
		if err := tpCmd.Run(ctx); err != nil {
			dieOnError("Failed to run start cmd", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
