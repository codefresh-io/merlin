package environment

import (
	"fmt"
	"io/ioutil"

	"github.com/codefresh-io/merlin/pkg/github"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/codefresh-io/merlin/pkg/template"

	"github.com/codefresh-io/merlin/pkg/codefresh"
	"github.com/codefresh-io/merlin/pkg/commander"
	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/kube"
	"github.com/codefresh-io/merlin/pkg/strvals"
)

type (
	Environment interface {
		Create(*CreateOptions) error
		Run(*RunCommandOptions) error
	}

	env struct {
		config *config.Config
		log    *logrus.Entry
	}

	RunCommandOptions struct {
		Component string
		Override  []string
		Command   string
	}

	CreateOptions struct {
		Name string
	}
)

func Build(c *config.Config, log *logrus.Entry) Environment {
	return &env{c, log}
}

func (e *env) Create(opt *CreateOptions) error {
	logger := e.log
	k, err := kube.New(e.config, logger)
	err = k.EnsureNamespaceNotExist(opt.Name)
	if err != nil {
		return err
	}
	logger.Debug("Namespace is not exist!")
	return codefresh.CreateEnvironment(&codefresh.Options{
		Name:   opt.Name,
		Config: e.config,
	}, logger)
}

func (e *env) readEnvironmentDescriptor() (*config.EnvironmentDescriptor, error) {
	if e.config.Environment.Path != "" {
		e.log.Debugf("Reading environment descriptor from: %s", e.config.Environment.Path)
		return config.ReadEnvironmentDescriptor(e.config.Environment.Path)
	} else {
		e.log.WithFields(map[string]interface{}{
			"Owner":    e.config.Environment.Git.Owner,
			"Repo":     e.config.Environment.Git.Repo,
			"Revision": e.config.Environment.Git.Revision,
		}).Debugf("Reading environment descriptor git , path: %s", e.config.Environment.Git.Path)
		g, err := github.New(e.config.Github.Token, e.log)
		if err != nil {
			return nil, err
		}
		git := e.config.Environment.Git
		content, err := g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
		if err != nil {
			return nil, err
		}
		descriptor := &config.EnvironmentDescriptor{}
		err = yaml.Unmarshal(content, descriptor)
		return descriptor, err
	}
}

func (e *env) readComponentTemplate(component *config.ComponentDescriptor) ([]byte, error) {
	if component.Spec.Path != "" {
		e.log.Debugf("Reading component template file from: %s", component.Spec.Path)
		return ioutil.ReadFile(component.Spec.Path)
	} else {
		e.log.WithFields(map[string]interface{}{
			"Owner":    component.Spec.Git.Owner,
			"Repo":     component.Spec.Git.Repo,
			"Revision": component.Spec.Git.Revision,
		}).Debugf("Reading component template from git : %s", component.Spec.Git.Path)
		g, err := github.New(e.config.Github.Token, e.log)
		if err != nil {
			return nil, err
		}
		git := component.Spec.Git
		content, err := g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
		if err != nil {
			return nil, err
		}
		return content, nil

	}
}

func (e *env) readValueFiles(component *config.ComponentDescriptor) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	for _, v := range component.Values {
		curr := map[string]interface{}{}
		if v.Path != "" {
			e.log.Debugf("Reading component template file from: %s", v.Path)
			content, err := ioutil.ReadFile(v.Path)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(content, &curr)
			if err != nil {
				return nil, err
			}
			mergeValues(base, curr)
		} else {
			e.log.WithFields(map[string]interface{}{
				"Owner":    v.Git.Owner,
				"Repo":     v.Git.Repo,
				"Revision": v.Git.Revision,
			}).Debugf("Reading component template from git : %s", v.Git.Path)
			g, err := github.New(e.config.Github.Token, e.log)
			if err != nil {
				return nil, err
			}
			git := v.Git
			content, err := g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(content, &curr)
			if err != nil {
				return nil, err
			}
			mergeValues(base, curr)
		}
	}
	return base, nil
}

func (e *env) Run(opt *RunCommandOptions) error {
	logger := e.log
	var component config.ComponentDescriptor
	rs := &config.RenderedService{}

	environmentPath := e.config.Environment.Path
	environmentDescriptor, err := e.readEnvironmentDescriptor()
	if err != nil {
		return err
	}
	for _, c := range environmentDescriptor.Components {
		if c.Name == opt.Component {
			component = c
		}
	}
	if &component == nil {
		return fmt.Errorf("Service: %s not found in Codefresh system config %s", opt.Component, environmentPath)
	}

	content, err := e.readComponentTemplate(&component)
	if err != nil {
		return err
	}

	// TODO read all files
	base, err := e.readValueFiles(&component)
	if err != nil {
		return err
	}

	dataSource := make(map[string]interface{})

	for _, value := range opt.Override {
		if err := strvals.ParseInto(value, base); err != nil {
			return fmt.Errorf("failed parsing --set data: %s", err)
		}
	}
	dataSource["Values"] = base

	system := make(map[string]interface{})
	jc, err := convertStruct(e.config)
	if err != nil {
		return err
	}
	system = mergeValues(system, jc)

	js, err := convertStruct(component)
	if err != nil {
		return err
	}
	system = mergeValues(system, js)

	dataSource["Merlin"] = system

	err = template.Render(content, dataSource, rs)
	if err != nil {
		return err
	}

	var cmd config.Command
	for _, command := range rs.Commands {

		if opt.Command == command.Name {
			logger.WithFields(map[string]interface{}{
				"Command": command.Name,
			}).Debug("Found command")
			cmd = command
		}

	}
	if &cmd != nil {
		logger.WithFields(map[string]interface{}{
			"Command": cmd.Name,
		}).Debug("Running step")
		logger.WithFields(map[string]interface{}{
			"Command": cmd.Name,
		}).Debugf("%s %v", cmd.Program, cmd.Args)
		err = commander.New(cmd.Program, cmd.Args, cmd.Env).Run()
		if err != nil {
			return err
		}
	}
	return err
}

func convertStruct(obj interface{}) (map[string]interface{}, error) {
	o, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	jo := make(map[string]interface{})
	err = yaml.Unmarshal(o, &jo)
	return jo, err
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
