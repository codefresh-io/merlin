package spec

import (
	"encoding/json"
	"fmt"
)

type (
	Component struct {
		Name string                 `json:"name"`
		Spec map[string]interface{} `json:"spec"`
	}
)

func (c *Component) ToJSObject() map[string]interface{} {
	res := map[string]interface{}{}
	b, err := json.Marshal(c)
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

func (c *Component) ToJSONString() string {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error marshalling: %s\n", err.Error())
		return ""
	}
	return string(b)
}
