package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/prismon/synthesis/internal/domain"
)

// TenantRepository handles tenant persistence
type TenantRepository struct {
	db *DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	labelsJSON, err := json.Marshal(tenant.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	query := `
		INSERT INTO tenant (id, owner, display_name, description, labels_json, version, last_modified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Owner,
		tenant.DisplayName,
		tenant.Description,
		labelsJSON,
		tenant.Version,
		tenant.LastModified,
	)

	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	// Update resource index
	if err := r.updateResourceIndex(ctx, tenant); err != nil {
		return fmt.Errorf("failed to update resource index: %w", err)
	}

	return nil
}

// Get retrieves a tenant by ID
func (r *TenantRepository) Get(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `
		SELECT id, owner, display_name, description, labels_json, version, last_modified
		FROM tenant
		WHERE id = $1
	`

	tenant := &domain.Tenant{}
	var labelsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Owner,
		&tenant.DisplayName,
		&tenant.Description,
		&labelsJSON,
		&tenant.Version,
		&tenant.LastModified,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if err := json.Unmarshal(labelsJSON, &tenant.Labels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
	}

	return tenant, nil
}

// List retrieves all tenants
func (r *TenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	query := `
		SELECT id, owner, display_name, description, labels_json, version, last_modified
		FROM tenant
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []*domain.Tenant

	for rows.Next() {
		tenant := &domain.Tenant{}
		var labelsJSON []byte

		err := rows.Scan(
			&tenant.ID,
			&tenant.Owner,
			&tenant.DisplayName,
			&tenant.Description,
			&labelsJSON,
			&tenant.Version,
			&tenant.LastModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}

		if err := json.Unmarshal(labelsJSON, &tenant.Labels); err != nil {
			return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
		}

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tenants: %w", err)
	}

	return tenants, nil
}

// Update updates an existing tenant
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	labelsJSON, err := json.Marshal(tenant.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	query := `
		UPDATE tenant
		SET owner = $2, display_name = $3, description = $4, labels_json = $5, version = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Owner,
		tenant.DisplayName,
		tenant.Description,
		labelsJSON,
		tenant.Version,
	)

	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("tenant not found: %s", tenant.ID)
	}

	// Update resource index
	if err := r.updateResourceIndex(ctx, tenant); err != nil {
		return fmt.Errorf("failed to update resource index: %w", err)
	}

	return nil
}

// Delete deletes a tenant
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tenant WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("tenant not found: %s", id)
	}

	// Delete from resource index
	_, err = r.db.ExecContext(ctx, `DELETE FROM resource_index WHERE entity_type = 'tenant' AND entity_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete from resource index: %w", err)
	}

	return nil
}

// updateResourceIndex updates the resource index for a tenant
func (r *TenantRepository) updateResourceIndex(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO resource_index (uri, entity_type, entity_id, tenant_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uri) DO UPDATE
		SET entity_type = EXCLUDED.entity_type,
			entity_id = EXCLUDED.entity_id,
			tenant_id = EXCLUDED.tenant_id
	`

	_, err := r.db.ExecContext(ctx, query, tenant.URI(), "tenant", tenant.ID, tenant.ID)
	return err
}
