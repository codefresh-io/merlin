package spec

import (
	"encoding/json"
	"fmt"
)

type (
	MerlinConfig struct {
		Name          string `json:"name"`
		EnvironmentJS string `json:"environment"`
		Codefresh     struct {
			Path    string `json:"path"`
			Context string `json:"context"`
		} `json:"codefresh"`
		Kubernetes struct {
			Path    string `json:"path"`
			Context string `json:"context"`
		} `json:"kubernetes"`
		Values map[string]interface{} `json:"values"`
	}
)

func (m *MerlinConfig) ToJSON() map[string]interface{} {
	res := map[string]interface{}{}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error marshalling: %s\n", err.Error())
		return res
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Printf("Error unmarshalling: %s\n", err.Error())
		return res
	}
	return res
}
