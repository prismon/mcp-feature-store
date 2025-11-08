package domain

// Product represents a business product associated with a tenant
type Product struct {
	TenantID    string `json:"tenantId" db:"tenant_id"`
	ID          string `json:"productId" db:"id"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
}

// ProductUser represents a user's association with a product
type ProductUser struct {
	ProductID string `json:"productId" db:"product_id"`
	UserID    string `json:"userId" db:"user_id"`
	Role      string `json:"role" db:"role"`
}

// URI returns the MCP URI for this product
func (p *Product) URI() string {
	return "synthesis://tenant/" + p.TenantID + "/product/" + p.ID
}
