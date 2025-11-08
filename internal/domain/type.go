package domain

// TypeDef represents content type definitions, renderers, editors, and constraints
type TypeDef struct {
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Renderers   []Renderer        `json:"renderers" db:"-"`
	Editors     []Editor          `json:"editors" db:"-"`
	Constraints []Constraint      `json:"constraints" db:"-"`
	Labels      map[string]string `json:"labels" db:"labels_json"`
}

// Renderer defines how a type is rendered
type Renderer struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// Editor defines how a type can be edited
type Editor struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// Constraint defines validation rules for a type
type Constraint struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// URI returns the MCP URI for this type definition
func (t *TypeDef) URI() string {
	return "synthesis://type/" + t.Name
}
