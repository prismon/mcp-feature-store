package domain

// Library represents a collection of notebooks within a tenant
type Library struct {
	TenantID    string            `json:"tenantId" db:"tenant_id"`
	ID          string            `json:"libraryId" db:"id"`
	Owner       string            `json:"owner" db:"owner"`
	DisplayName string            `json:"display_name" db:"display_name"`
	Description string            `json:"description" db:"description"`
	Labels      map[string]string `json:"labels" db:"labels_json"`
}

// URI returns the MCP URI for this library
func (l *Library) URI() string {
	return "synthesis://tenant/" + l.TenantID + "/library/" + l.ID
}
