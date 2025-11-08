package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prismon/synthesis/internal/domain"
)

// NotebookRepository handles notebook persistence
type NotebookRepository struct {
	db *DB
}

// NewNotebookRepository creates a new notebook repository
func NewNotebookRepository(db *DB) *NotebookRepository {
	return &NotebookRepository{db: db}
}

// Create creates a new notebook with its content
func (r *NotebookRepository) Create(ctx context.Context, notebook *domain.Notebook) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert notebook
	query := `
		INSERT INTO notebook (id, tenant_id, library_id, status, owner, display_name, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(ctx, query,
		notebook.ID,
		notebook.TenantID,
		notebook.LibraryID,
		notebook.Status,
		notebook.Owner,
		notebook.DisplayName,
		notebook.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create notebook: %w", err)
	}

	// Insert notebook content
	contentQuery := `
		INSERT INTO notebook_content (notebook_id, markdown)
		VALUES ($1, $2)
	`

	_, err = tx.ExecContext(ctx, contentQuery, notebook.ID, notebook.Contents.Data.Markdown)
	if err != nil {
		return fmt.Errorf("failed to create notebook content: %w", err)
	}

	// Insert content blocks
	for _, block := range notebook.Contents.ContentBlocks {
		if err := r.createContentBlock(ctx, tx, notebook.ID, &block); err != nil {
			return fmt.Errorf("failed to create content block: %w", err)
		}
	}

	// Insert notifications
	for _, notif := range notebook.Notifications {
		notifQuery := `INSERT INTO notebook_notification (notebook_id, nurl) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, notifQuery, notebook.ID, notif.URL)
		if err != nil {
			return fmt.Errorf("failed to create notification: %w", err)
		}
	}

	// Update resource index
	if err := r.updateResourceIndex(ctx, tx, notebook); err != nil {
		return fmt.Errorf("failed to update resource index: %w", err)
	}

	return tx.Commit()
}

// Get retrieves a notebook by ID
func (r *NotebookRepository) Get(ctx context.Context, id string) (*domain.Notebook, error) {
	query := `
		SELECT id, tenant_id, library_id, status, owner, display_name, description
		FROM notebook
		WHERE id = $1
	`

	notebook := &domain.Notebook{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notebook.ID,
		&notebook.TenantID,
		&notebook.LibraryID,
		&notebook.Status,
		&notebook.Owner,
		&notebook.DisplayName,
		&notebook.Description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notebook not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notebook: %w", err)
	}

	// Get content
	contentQuery := `SELECT markdown FROM notebook_content WHERE notebook_id = $1`
	err = r.db.QueryRowContext(ctx, contentQuery, id).Scan(&notebook.Contents.Data.Markdown)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get notebook content: %w", err)
	}

	// Get content blocks
	blocks, err := r.getContentBlocks(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content blocks: %w", err)
	}
	notebook.Contents.ContentBlocks = blocks

	// Get notifications
	notifs, err := r.getNotifications(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	notebook.Notifications = notifs

	return notebook, nil
}

// ListByLibrary retrieves all notebooks in a library
func (r *NotebookRepository) ListByLibrary(ctx context.Context, libraryID string) ([]*domain.Notebook, error) {
	query := `
		SELECT id, tenant_id, library_id, status, owner, display_name, description
		FROM notebook
		WHERE library_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, libraryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notebooks: %w", err)
	}
	defer rows.Close()

	var notebooks []*domain.Notebook

	for rows.Next() {
		notebook := &domain.Notebook{}
		err := rows.Scan(
			&notebook.ID,
			&notebook.TenantID,
			&notebook.LibraryID,
			&notebook.Status,
			&notebook.Owner,
			&notebook.DisplayName,
			&notebook.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notebook: %w", err)
		}

		notebooks = append(notebooks, notebook)
	}

	return notebooks, rows.Err()
}

// Update updates an existing notebook
func (r *NotebookRepository) Update(ctx context.Context, notebook *domain.Notebook) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE notebook
		SET status = $2, owner = $3, display_name = $4, description = $5
		WHERE id = $1
	`

	result, err := tx.ExecContext(ctx, query,
		notebook.ID,
		notebook.Status,
		notebook.Owner,
		notebook.DisplayName,
		notebook.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to update notebook: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("notebook not found: %s", notebook.ID)
	}

	// Update content
	contentQuery := `
		UPDATE notebook_content
		SET markdown = $2
		WHERE notebook_id = $1
	`

	_, err = tx.ExecContext(ctx, contentQuery, notebook.ID, notebook.Contents.Data.Markdown)
	if err != nil {
		return fmt.Errorf("failed to update notebook content: %w", err)
	}

	return tx.Commit()
}

// Delete deletes a notebook
func (r *NotebookRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notebook WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("notebook not found: %s", id)
	}

	return nil
}

// Helper methods

func (r *NotebookRepository) createContentBlock(ctx context.Context, tx *sql.Tx, notebookID string, block *domain.ContentBlock) error {
	query := `
		INSERT INTO content_block (notebook_id, uid, parent_uid, content_type, data, "order")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var blockID string
	err := tx.QueryRowContext(ctx, query,
		notebookID,
		block.UID,
		block.ParentUID,
		block.ContentType,
		block.Data,
		block.Order,
	).Scan(&blockID)

	if err != nil {
		return err
	}

	// Insert block types
	for _, typeName := range block.Types {
		typeQuery := `INSERT INTO content_block_type (content_block_id, type_name) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, typeQuery, blockID, typeName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *NotebookRepository) getContentBlocks(ctx context.Context, notebookID string) ([]domain.ContentBlock, error) {
	query := `
		SELECT id, uid, parent_uid, content_type, data, "order"
		FROM content_block
		WHERE notebook_id = $1
		ORDER BY "order"
	`

	rows, err := r.db.QueryContext(ctx, query, notebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []domain.ContentBlock

	for rows.Next() {
		var block domain.ContentBlock
		var blockID string

		err := rows.Scan(
			&blockID,
			&block.UID,
			&block.ParentUID,
			&block.ContentType,
			&block.Data,
			&block.Order,
		)
		if err != nil {
			return nil, err
		}

		// Get types for this block
		typeQuery := `SELECT type_name FROM content_block_type WHERE content_block_id = $1`
		typeRows, err := r.db.QueryContext(ctx, typeQuery, blockID)
		if err != nil {
			return nil, err
		}

		var types []string
		for typeRows.Next() {
			var typeName string
			if err := typeRows.Scan(&typeName); err != nil {
				typeRows.Close()
				return nil, err
			}
			types = append(types, typeName)
		}
		typeRows.Close()

		block.Types = types
		blocks = append(blocks, block)
	}

	return blocks, rows.Err()
}

func (r *NotebookRepository) getNotifications(ctx context.Context, notebookID string) ([]domain.Notification, error) {
	query := `SELECT nurl FROM notebook_notification WHERE notebook_id = $1`

	rows, err := r.db.QueryContext(ctx, query, notebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification

	for rows.Next() {
		var notif domain.Notification
		if err := rows.Scan(&notif.URL); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, rows.Err()
}

func (r *NotebookRepository) updateResourceIndex(ctx context.Context, tx *sql.Tx, notebook *domain.Notebook) error {
	query := `
		INSERT INTO resource_index (uri, entity_type, entity_id, tenant_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uri) DO UPDATE
		SET entity_type = EXCLUDED.entity_type,
			entity_id = EXCLUDED.entity_id,
			tenant_id = EXCLUDED.tenant_id
	`

	_, err := tx.ExecContext(ctx, query, notebook.URI(), "notebook", notebook.ID, notebook.TenantID)
	return err
}
