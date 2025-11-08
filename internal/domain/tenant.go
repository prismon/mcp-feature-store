package domain

import "time"

// Tenant represents a top-level organization unit
type Tenant struct {
	ID           string            `json:"tenantId" db:"id"`
	Owner        string            `json:"owner" db:"owner"`
	DisplayName  string            `json:"display_name" db:"display_name"`
	Description  string            `json:"description" db:"description"`
	Labels       map[string]string `json:"labels" db:"labels_json"`
	Version      string            `json:"version" db:"version"`
	LastModified time.Time         `json:"last_modified" db:"last_modified"`
}

// URI returns the MCP URI for this tenant
func (t *Tenant) URI() string {
	return "synthesis://tenant/" + t.ID
}
