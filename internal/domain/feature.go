package domain

import "time"

// Feature represents derived data associated with resources
type Feature struct {
	TenantID      string              `json:"tenantId" db:"tenant_id"`
	ID            string              `json:"featureId" db:"id"`
	DisplayName   string              `json:"display_name" db:"display_name"`
	Description   string              `json:"description" db:"description"`
	Resources     []ExternalResource  `json:"resources,omitempty"`
	Notifications []Notification      `json:"notification,omitempty"`
	TTL           time.Duration       `json:"ttl,omitempty" db:"ttl"`
	Values        map[string]string   `json:"values" db:"values_json"`
}

// ExternalResource represents a URL reference to an external resource
type ExternalResource struct {
	URL string `json:"url" db:"url"`
}

// URI returns the MCP URI for this feature
func (f *Feature) URI() string {
	return "synthesis://tenant/" + f.TenantID + "/feature/" + f.ID
}
