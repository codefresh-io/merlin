package template

import (
	"bytes"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/codefresh-io/merlin/pkg/utils"
	"github.com/hairyhenderson/gomplate"
	"github.com/imdario/mergo"
)

func Render(toRender []byte, data interface{}, target interface{}) error {
	out := new(bytes.Buffer)
	tmpl := template.New("")
	funcMap, err := getFuncs()
	if err != nil {
		return err
	}
	tmpl.Funcs(funcMap)
	tmpl.Parse(string(toRender))
	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(out.Bytes(), target)
	return err
}

func getFuncs() (template.FuncMap, error) {
	funcs := map[string]interface{}{
		"GetAvailablePort": utils.GetAvailablePort,
	}
	if err := mergo.Map(&funcs, gomplate.Funcs(nil)); err != nil {
		return nil, err
	}
	return funcs, nil
}
