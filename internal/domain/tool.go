package domain

// ToolConfig represents configuration for external tools/integrations
type ToolConfig struct {
	TenantID    string                 `json:"tenantId" db:"tenant_id"`
	ID          string                 `json:"toolId" db:"id"`
	DisplayName string                 `json:"display_name" db:"display_name"`
	Description string                 `json:"description" db:"description"`
	Config      map[string]interface{} `json:"config" db:"config_json"`
}

// URI returns the MCP URI for this tool configuration
func (t *ToolConfig) URI() string {
	return "synthesis://tenant/" + t.TenantID + "/tool/" + t.ID
}
