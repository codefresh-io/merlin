package environment

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codefresh-io/merlin/pkg/github"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/codefresh-io/merlin/pkg/template"

	"github.com/codefresh-io/merlin/pkg/cache"
	"github.com/codefresh-io/merlin/pkg/commander"
	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/codefresh-io/merlin/pkg/strvals"
)

type (
	Environment interface {
		Run(*RunCommandOptions) error
		List(*ListCommandOptions) ([][]string, error)
	}

	env struct {
		config *config.Config
		log    *logrus.Entry
		cache  cache.Cache
	}

	signalHandler interface {
		Process()
		Push(chan os.Signal)
	}

	RunCommandOptions struct {
		Component     string
		Override      []string
		Operator      string
		SkipExec      bool
		SignalHandler signalHandler
	}

	ListCommandOptions struct{}
)

func Build(c *config.Config, cache cache.Cache, log *logrus.Entry) Environment {
	return &env{c, log, cache}
}

func (e *env) readComponentTemplate(component *config.Component) ([]byte, error) {
	content := []byte{}
	if component.Spec.Path != "" {
		e.log.Debugf("Reading component template file from: %s", component.Spec.Path)
		return ioutil.ReadFile(component.Spec.Path)
	} else {
		git := component.Spec.Git
		key := fmt.Sprintf("%s.%s.%s.%s", git.Owner, git.Repo, git.Path, git.Revision)
		e.log.WithFields(map[string]interface{}{
			"Owner":    git.Owner,
			"Repo":     git.Repo,
			"Revision": git.Revision,
		}).Debugf("Reading component template from git : %s", git.Path)
		if err := e.cache.Get(key, &content); err != nil {
			g := github.New(e.config.Github.Token, e.log)
			content, err = g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
			if err != nil {
				return nil, err
			}
			e.log.Debug("Saving component template to cache")
			e.cache.Put(key, &content)
			return content, nil
		}
		e.log.Debug("Component template loaded from cache")
		return content, nil

	}
}

func (e *env) readEnvironmentTemplate() ([]byte, error) {
	content := []byte{}
	pathToLocalEnvironment := e.config.Environment.Spec.Path
	if pathToLocalEnvironment != "" {
		e.log.Debugf("Reading environment template file from: %s", pathToLocalEnvironment)
		return ioutil.ReadFile(pathToLocalEnvironment)
	} else {
		git := e.config.Environment.Spec.Git
		key := fmt.Sprintf("%s.%s.%s.%s", git.Owner, git.Repo, git.Path, git.Revision)
		e.log.WithFields(map[string]interface{}{
			"Owner":    git.Owner,
			"Repo":     git.Repo,
			"Revision": git.Revision,
		}).Debugf("Reading environment template from git : %s", git.Path)
		if err := e.cache.Get(key, &content); err != nil {
			g := github.New(e.config.Github.Token, e.log)
			content, err = g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
			if err != nil {
				return nil, err
			}
			e.log.Debug("Saving environment template to cache")
			e.cache.Put(key, &content)
			return content, nil
		}
		e.log.Debug("Environment template loaded from cache")
		return content, nil
	}
}

func (e *env) readValueFiles(values []config.Values, override []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	for _, v := range values {
		curr := map[string]interface{}{}
		if v.Path != "" {
			e.log.Debugf("Reading value files from: %s", v.Path)
			content, err := ioutil.ReadFile(v.Path)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(content, &curr)
			if err != nil {
				return nil, err
			}
			base = mergeValues(base, curr)
		} else {
			git := v.Git
			key := fmt.Sprintf("%s.%s.%s.%s", git.Owner, git.Repo, git.Path, git.Revision)
			e.log.WithFields(map[string]interface{}{
				"Owner":    git.Owner,
				"Repo":     git.Repo,
				"Revision": git.Revision,
			}).Debugf("Reading value files from git : %s", git.Path)
			if err := e.cache.Get(key, &curr); err != nil {
				g := github.New(e.config.Github.Token, e.log)
				content, err := g.ReadFile(git.Owner, git.Repo, git.Path, git.Revision)
				if err != nil {
					return nil, err
				}
				err = yaml.Unmarshal(content, &curr)
				if err != nil {
					return nil, err
				}
				e.log.Debug("Saving component values to cache")
				e.cache.Put(key, &curr)
			} else {
				e.log.Debug("Component values loaded from cache")
			}
			base = mergeValues(base, curr)
		}
	}

	for _, value := range override {
		e.log.Debugf("Recieved overwrite value: %s", value)
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return base, nil
}

func (e *env) prepareEnvironmemtDescriptor(override []string) (*config.EnvironmentDescriptor, error) {
	source := make(map[string]interface{})
	environmentDescriptor := &config.EnvironmentDescriptor{}
	environmentContent, err := e.readEnvironmentTemplate()
	if err != nil {
		return nil, err
	}

	systemSet, err := e.getSystemVariables()
	if err != nil {
		return nil, err
	}
	source["Merlin"] = systemSet

	values, err := e.readValueFiles(e.config.Environment.Values, override)
	if err != nil {
		return nil, err
	}
	source["Values"] = values

	err = template.Render(environmentContent, source, environmentDescriptor)
	return environmentDescriptor, err
}

func (e *env) getSystemVariables() (map[string]interface{}, error) {
	system := make(map[string]interface{})
	json, err := convertStruct(e.config)
	if err != nil {
		return nil, err
	}
	return mergeValues(system, json), nil
}

func (e *env) prepareComponentDescriptor(override []string, component *config.Component) (*config.ComponentDescriptor, error) {
	source := make(map[string]interface{})
	componentDescriptor := &config.ComponentDescriptor{}

	componentContent, err := e.readComponentTemplate(component)
	if err != nil {
		return nil, err
	}

	systemSet, err := e.getSystemVariables()
	if err != nil {
		return nil, err
	}
	source["Merlin"] = systemSet

	values, err := e.readValueFiles(component.Values, override)
	if err != nil {
		return nil, err
	}
	source["Values"] = values

	componentValues := make(map[string]interface{})
	componentStruct, err := convertStruct(component)
	if err != nil {
		return nil, err
	}
	source["Component"] = mergeValues(componentValues, componentStruct)

	err = template.Render(componentContent, source, componentDescriptor)
	if err != nil {
		return nil, err
	}
	return componentDescriptor, err
}

func (e *env) getCandidateOperators() {

}

func (e *env) Run(opt *RunCommandOptions) error {
	var component config.Component
	operators := []config.Operator{}
	logger := e.log

	environmentDescriptor, err := e.prepareEnvironmemtDescriptor(opt.Override)
	if err != nil {
		return err
	}
	for _, o := range environmentDescriptor.Operators {
		if o.Name == opt.Operator {
			logger.Debugf("Adding operator %s to operator slice", o.Name)
			operators = append(operators, o)
		}
	}

	// setup envs to be added to the operator
	envs := []string{}

	if opt.Component != "" {
		found := false
		for _, c := range environmentDescriptor.Components {
			if c.Name == opt.Component {
				component = c
				found = true
			}
		}
		if found {

			logger.Debugf("Setting MERLIN_COMPONENT=%s", component.Name)
			envs = append(envs, fmt.Sprintf("MERLIN_COMPONENT=%s", component.Name))
			componentDescriptor, err := e.prepareComponentDescriptor(opt.Override, &component)
			if err != nil {
				return err
			}
			for _, op := range componentDescriptor.Operators {

				if opt.Operator == op.Name {
					logger.Debugf("Adding operator %s to operator slice", op.Name)
					operators = append(operators, op)
				}
			}
		} else {
			return fmt.Errorf("Component %s not found", opt.Component)
		}
	}

	for _, o := range operators {
		var operatorEnv []string
		if o.Description != "" {
			logger.Infof("Running: %s", o.Description)
		}

		operatorEnv = append(envs, o.Spec.Env...)
		if opt.SkipExec {
			logger.Debugf("Skipping execution, actual command:\n%v", append([]string{o.Spec.Program}, o.Spec.Args...))
		} else {
			err = commander.New(&commander.Options{
				Program:       o.Spec.Program,
				Args:          o.Spec.Args,
				Env:           operatorEnv,
				SignalHandler: opt.SignalHandler,
				Detached:      o.Spec.Detached,
				Logger:        logger,
			}).Run()
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (e *env) List(_ *ListCommandOptions) ([][]string, error) {
	res := [][]string{}
	environmentDescriptor, err := e.prepareEnvironmemtDescriptor(nil)
	if err != nil {
		return nil, err
	}
	for _, op := range environmentDescriptor.Operators {
		row := []string{"Envronment", op.Name, op.Description}
		res = append(res, row)
	}

	for _, c := range environmentDescriptor.Components {
		e.log.Debugf("Getting info about component %s", c.Name)
		componentDescriptor, err := e.prepareComponentDescriptor(nil, &c)
		if err != nil {
			return nil, err
		}
		for _, op := range componentDescriptor.Operators {
			row := []string{fmt.Sprintf("Component: [%s]", c.Name), op.Name, op.Description}
			res = append(res, row)
		}
	}
	return res, nil
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
