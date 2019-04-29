package js

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/robertkrimen/otto"

	// Add underscore.js to otto
	"github.com/codefresh-io/merlin/pkg/utils"
	_ "github.com/robertkrimen/otto/underscore"
)

const (
	undefinedError = "undefined"
)

type (
	loader interface {
	}
	functionCaller interface {
	}

	Runner interface {
		Load(scripts []string, gloabl map[string]interface{}) (*otto.Value, error)
		CallFn(name string, input ...interface{}) (*otto.Value, error)
	}

	runner struct {
		vm *otto.Otto
	}
)

func NewJSEngine() Runner {
	return &runner{
		vm: otto.New(),
	}
}

func (r *runner) Load(scripts []string, input map[string]interface{}) (*otto.Value, error) {
	global := TemplatesMap()["classes.js"]
	for _, script := range scripts {
		c, err := r.vm.Compile("", string(script))
		if err != nil {
			return nil, err
		}
		global = fmt.Sprintf("%s\n%s", global, c.String())
	}
	for k, v := range input {
		err := r.vm.Set(k, v)
		if err != nil {
			return nil, err
		}
	}
	r.vm.Set("GetAvailablePort", func(call otto.FunctionCall) otto.Value {
		p, err := utils.GetAvailablePort()
		if err != nil {
			v, _ := r.vm.ToValue(err.Error())
			return v
		}
		v, _ := r.vm.ToValue(p)
		return v
	})
	r.vm.Set("process", map[string]interface{}{
		"env": getProcessEnv(),
	})
	res, err := r.vm.Run(global)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *runner) CallFn(name string, input ...interface{}) (*otto.Value, error) {
	v, err := r.vm.Get(name)
	if err != nil {
		return nil, err
	}
	if v.IsDefined() {
		res, err := v.Call(otto.Value{}, input...)
		return &res, err
	}
	return nil, errors.New(undefinedError)
}

func getProcessEnv() map[string]string {
	processEnv := map[string]string{}
	{
		for _, env := range os.Environ() {
			res := strings.Split(env, "=")
			processEnv[res[0]] = strings.Join(res[1:], "=")
		}
	}
	return processEnv
}
