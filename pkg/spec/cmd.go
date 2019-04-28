package spec

type (
	CmdSet []Cmd

	Cmd struct {
		Name     string   `json:"name"`
		Exec     []string `json:"exec"`
		Env      []string `json:"env"`
		Detached bool     `json:"detached"`
		WorkDir  string   `json:"workDir"`
	}
)
