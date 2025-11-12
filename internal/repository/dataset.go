package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/supercakecrumb/affirm-name/internal/model"
)

// DatasetRepository handles database operations for datasets
type DatasetRepository struct {
	db *sqlx.DB
}

// NewDatasetRepository creates a new DatasetRepository
func NewDatasetRepository(db *sqlx.DB) *DatasetRepository {
	return &DatasetRepository{db: db}
}

// Create creates a new dataset record
func (r *DatasetRepository) Create(ctx context.Context, req *model.CreateDatasetRequest) (*model.Dataset, error) {
	dataset := &model.Dataset{
		ID:         uuid.New(),
		CountryID:  req.CountryID,
		Filename:   req.Filename,
		FilePath:   req.FilePath,
		FileSize:   req.FileSize,
		Status:     model.DatasetStatusPending,
		UploadedBy: req.UploadedBy,
		UploadedAt: time.Now().UTC(),
	}

	query := `
		INSERT INTO datasets (
			id, country_id, filename, file_path, file_size, 
			status, uploaded_by, uploaded_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		dataset.ID,
		dataset.CountryID,
		dataset.Filename,
		dataset.FilePath,
		dataset.FileSize,
		dataset.Status,
		dataset.UploadedBy,
		dataset.UploadedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataset: %w", err)
	}

	return dataset, nil
}

// GetByID retrieves a dataset by ID
func (r *DatasetRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Dataset, error) {
	var dataset model.Dataset

	query := `
		SELECT 
			id, country_id, filename, file_path, file_size,
			status, row_count, error_message, uploaded_by,
			uploaded_at, processed_at, deleted_at
		FROM datasets
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &dataset, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dataset not found")
		}
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}

	return &dataset, nil
}

// GetByIDWithCountry retrieves a dataset with country information
func (r *DatasetRepository) GetByIDWithCountry(ctx context.Context, id uuid.UUID) (*model.DatasetWithCountry, error) {
	var dataset model.DatasetWithCountry

	query := `
		SELECT 
			d.id, d.country_id, d.filename, d.file_path, d.file_size,
			d.status, d.row_count, d.error_message, d.uploaded_by,
			d.uploaded_at, d.processed_at, d.deleted_at,
			c.code as country_code, c.name as country_name
		FROM datasets d
		JOIN countries c ON c.id = d.country_id
		WHERE d.id = $1 AND d.deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &dataset, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dataset not found")
		}
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}

	return &dataset, nil
}

// List retrieves datasets with optional filters
func (r *DatasetRepository) List(ctx context.Context, filters *model.DatasetFilters, limit, offset int) ([]*model.DatasetWithCountry, int, error) {
	// Build query with filters
	query := `
		SELECT 
			d.id, d.country_id, d.filename, d.file_path, d.file_size,
			d.status, d.row_count, d.error_message, d.uploaded_by,
			d.uploaded_at, d.processed_at, d.deleted_at,
			c.code as country_code, c.name as country_name
		FROM datasets d
		JOIN countries c ON c.id = d.country_id
		WHERE d.deleted_at IS NULL
	`

	countQuery := `
		SELECT COUNT(*)
		FROM datasets d
		WHERE d.deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 1

	// Apply filters
	if filters != nil {
		if filters.CountryID != nil {
			query += fmt.Sprintf(" AND d.country_id = $%d", argCount)
			countQuery += fmt.Sprintf(" AND d.country_id = $%d", argCount)
			args = append(args, *filters.CountryID)
			argCount++
		}
		if filters.Status != nil {
			query += fmt.Sprintf(" AND d.status = $%d", argCount)
			countQuery += fmt.Sprintf(" AND d.status = $%d", argCount)
			args = append(args, *filters.Status)
			argCount++
		}
	}

	// Get total count
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count datasets: %w", err)
	}

	// Add ordering and pagination
	query += " ORDER BY d.uploaded_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Execute query
	var datasets []*model.DatasetWithCountry
	err = r.db.SelectContext(ctx, &datasets, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list datasets: %w", err)
	}

	return datasets, total, nil
}

// UpdateStatus updates the status of a dataset
func (r *DatasetRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.DatasetStatus) error {
	query := `
		UPDATE datasets
		SET status = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("dataset not found")
	}

	return nil
}

// UpdateCompleted marks a dataset as completed with row count
func (r *DatasetRepository) UpdateCompleted(ctx context.Context, id uuid.UUID, rowCount int) error {
	now := time.Now().UTC()

	query := `
		UPDATE datasets
		SET 
			status = $1,
			row_count = $2,
			processed_at = $3,
			error_message = NULL
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, model.DatasetStatusCompleted, rowCount, now, id)
	if err != nil {
		return fmt.Errorf("failed to update dataset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("dataset not found")
	}

	return nil
}

// UpdateFailed marks a dataset as failed with error message
func (r *DatasetRepository) UpdateFailed(ctx context.Context, id uuid.UUID, errorMessage string) error {
	now := time.Now().UTC()

	query := `
		UPDATE datasets
		SET 
			status = $1,
			error_message = $2,
			processed_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, model.DatasetStatusFailed, errorMessage, now, id)
	if err != nil {
		return fmt.Errorf("failed to update dataset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("dataset not found")
	}

	return nil
}

// Delete soft deletes a dataset
func (r *DatasetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()

	query := `
		UPDATE datasets
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete dataset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("dataset not found")
	}

	return nil
}
