package codefresh

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/codefresh-io/merlin/pkg/config"

	sdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	utils "github.com/codefresh-io/go-sdk/pkg/utils"
)

type (
	Options struct {
		Name   string
		Config *config.Config
	}
)

func CreateEnvironment(opt *Options, log *logrus.Entry) error {
	cf, err := getCodefreshClient(opt.Config, log)
	if err != nil {
		return err
	}
	id, err := cf.Pipelines().Run("codefresh-io/cf-helm/create-dynamic-env", &sdk.RunOptions{
		Branch: "master",
		Variables: map[string]string{
			"NAMESPACE_NAME": opt.Name,
			"WAIT":           "false",
		},
	})
	if err != nil {
		return err
	}
	log.WithField("Workflow", id).Debug("Creating environment by running pipeline in Codefresh")
	err = cf.Workflows().WaitForStatus(id, "success", 5*time.Second, 10*time.Minute)
	if err != nil {
		return err
	}
	log.Debug("Created!")
	return nil
}

func getCodefreshClient(c *config.Config, log *logrus.Entry) (sdk.Codefresh, error) {
	log.Debugf("Reading context %s from %s", c.Codefresh.Path, c.Codefresh.Context)
	context, err := utils.ReadAuthContext(c.Codefresh.Path, c.Codefresh.Context)
	if err != nil {
		return nil, err
	}
	clientOptions := sdk.ClientOptions{Host: context.URL,
		Auth: sdk.AuthOptions{Token: context.Token}}
	cf := sdk.New(&clientOptions)
	return cf, nil
}
