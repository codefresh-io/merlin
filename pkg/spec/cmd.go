package spec

type (
	CmdSet []Cmd

	Cmd struct {
		Name     string   `json:"name"`
		Exec     []string `json:"exec"`
		Program  string   `json:"program"`
		Env      []string `json:"env"`
		Detached bool     `json:"detached"`
		WorkDir  string   `json:"workDir"`
	}
)
