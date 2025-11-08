package postgres

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
)

// VectorRepository handles vector operations for semantic search
type VectorRepository struct {
	db *DB
}

// NewVectorRepository creates a new vector repository
func NewVectorRepository(db *DB) *VectorRepository {
	return &VectorRepository{db: db}
}

// NotebookEmbedding represents a notebook with its embedding
type NotebookEmbedding struct {
	NotebookID string
	Embedding  pgvector.Vector
	Model      string
}

// FeatureEmbedding represents a feature with its embedding
type FeatureEmbedding struct {
	FeatureID string
	Embedding pgvector.Vector
	Model     string
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	ID         string
	Similarity float64
}

// UpsertNotebookEmbedding inserts or updates a notebook embedding
func (r *VectorRepository) UpsertNotebookEmbedding(ctx context.Context, notebookID string, embedding []float32, model string) error {
	query := `
		INSERT INTO notebook_embedding (notebook_id, embedding, model)
		VALUES ($1, $2, $3)
		ON CONFLICT (notebook_id) DO UPDATE
		SET embedding = EXCLUDED.embedding,
			model = EXCLUDED.model,
			updated_at = CURRENT_TIMESTAMP
	`

	vec := pgvector.NewVector(embedding)

	_, err := r.db.ExecContext(ctx, query, notebookID, vec, model)
	if err != nil {
		return fmt.Errorf("failed to upsert notebook embedding: %w", err)
	}

	return nil
}

// UpsertFeatureEmbedding inserts or updates a feature embedding
func (r *VectorRepository) UpsertFeatureEmbedding(ctx context.Context, featureID string, embedding []float32, model string) error {
	query := `
		INSERT INTO feature_embedding (feature_id, embedding, model)
		VALUES ($1, $2, $3)
		ON CONFLICT (feature_id) DO UPDATE
		SET embedding = EXCLUDED.embedding,
			model = EXCLUDED.model,
			updated_at = CURRENT_TIMESTAMP
	`

	vec := pgvector.NewVector(embedding)

	_, err := r.db.ExecContext(ctx, query, featureID, vec, model)
	if err != nil {
		return fmt.Errorf("failed to upsert feature embedding: %w", err)
	}

	return nil
}

// SearchNotebooks performs semantic search on notebooks using cosine similarity
func (r *VectorRepository) SearchNotebooks(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	query := `
		SELECT notebook_id, 1 - (embedding <=> $1) as similarity
		FROM notebook_embedding
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	vec := pgvector.NewVector(queryEmbedding)

	rows, err := r.db.QueryContext(ctx, query, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search notebooks: %w", err)
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var result SearchResult
		err := rows.Scan(&result.ID, &result.Similarity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// SearchFeatures performs semantic search on features using cosine similarity
func (r *VectorRepository) SearchFeatures(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	query := `
		SELECT feature_id, 1 - (embedding <=> $1) as similarity
		FROM feature_embedding
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	vec := pgvector.NewVector(queryEmbedding)

	rows, err := r.db.QueryContext(ctx, query, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search features: %w", err)
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var result SearchResult
		err := rows.Scan(&result.ID, &result.Similarity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// SearchNotebooksByTenant performs semantic search on notebooks within a specific tenant
func (r *VectorRepository) SearchNotebooksByTenant(ctx context.Context, tenantID string, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	query := `
		SELECT ne.notebook_id, 1 - (ne.embedding <=> $1) as similarity
		FROM notebook_embedding ne
		JOIN notebook n ON ne.notebook_id = n.id
		WHERE n.tenant_id = $2
		ORDER BY ne.embedding <=> $1
		LIMIT $3
	`

	vec := pgvector.NewVector(queryEmbedding)

	rows, err := r.db.QueryContext(ctx, query, vec, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search notebooks by tenant: %w", err)
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var result SearchResult
		err := rows.Scan(&result.ID, &result.Similarity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// DeleteNotebookEmbedding deletes a notebook embedding
func (r *VectorRepository) DeleteNotebookEmbedding(ctx context.Context, notebookID string) error {
	query := `DELETE FROM notebook_embedding WHERE notebook_id = $1`

	_, err := r.db.ExecContext(ctx, query, notebookID)
	if err != nil {
		return fmt.Errorf("failed to delete notebook embedding: %w", err)
	}

	return nil
}

// DeleteFeatureEmbedding deletes a feature embedding
func (r *VectorRepository) DeleteFeatureEmbedding(ctx context.Context, featureID string) error {
	query := `DELETE FROM feature_embedding WHERE feature_id = $1`

	_, err := r.db.ExecContext(ctx, query, featureID)
	if err != nil {
		return fmt.Errorf("failed to delete feature embedding: %w", err)
	}

	return nil
}
