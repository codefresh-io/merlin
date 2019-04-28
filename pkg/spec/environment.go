package spec

type (
	Environment struct {
		Version    string      `json:"version"`
		Operators  []Operator  `json:"operators"`
		Components []Component `json:"components"`
	}
)
