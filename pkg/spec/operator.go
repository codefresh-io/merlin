package spec

var ScopeComponent = "component"

type (
	Operator struct {
		Name        string  `json:"name"`
		Params      []Param `json:"params"`
		Description string  `json:"description"`
		Scope       string  `json:"scope"`
	}

	Param struct {
		Name                string `json:"name"`
		Description         string `json:"description"`
		EnvironmentVariable string `json:"envVar"`
		Default             string `json:"default"`
		Required            bool   `json:"required"`
		InteractiveInput    bool   `json:"interactive"`
	}
)
