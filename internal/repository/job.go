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

// JobRepository handles database operations for jobs
type JobRepository struct {
	db *sqlx.DB
}

// NewJobRepository creates a new JobRepository
func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create creates a new job record
func (r *JobRepository) Create(ctx context.Context, req *model.CreateJobRequest) (*model.Job, error) {
	maxAttempts := req.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = 3 // Default
	}

	job := &model.Job{
		ID:          uuid.New(),
		DatasetID:   req.DatasetID,
		Type:        req.Type,
		Status:      model.JobStatusQueued,
		Payload:     req.Payload,
		Attempts:    0,
		MaxAttempts: maxAttempts,
		CreatedAt:   time.Now().UTC(),
	}

	query := `
		INSERT INTO jobs (
			id, dataset_id, type, status, payload, 
			attempts, max_attempts, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		job.ID,
		job.DatasetID,
		job.Type,
		job.Status,
		job.Payload,
		job.Attempts,
		job.MaxAttempts,
		job.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return job, nil
}

// GetByID retrieves a job by ID
func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	var job model.Job

	query := `
		SELECT 
			id, dataset_id, type, status, payload,
			attempts, max_attempts, last_error, next_retry_at,
			locked_at, locked_by, created_at, started_at, completed_at
		FROM jobs
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &job, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}

// List retrieves jobs with optional filters
func (r *JobRepository) List(ctx context.Context, filters *model.JobFilters, limit, offset int) ([]*model.Job, int, error) {
	// Build query with filters
	query := `
		SELECT 
			id, dataset_id, type, status, payload,
			attempts, max_attempts, last_error, next_retry_at,
			locked_at, locked_by, created_at, started_at, completed_at
		FROM jobs
		WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM jobs
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Apply filters
	if filters != nil {
		if filters.Status != nil {
			query += fmt.Sprintf(" AND status = $%d", argCount)
			countQuery += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, *filters.Status)
			argCount++
		}
		if filters.Type != nil {
			query += fmt.Sprintf(" AND type = $%d", argCount)
			countQuery += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, *filters.Type)
			argCount++
		}
		if filters.DatasetID != nil {
			query += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			countQuery += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			args = append(args, *filters.DatasetID)
			argCount++
		}
	}

	// Get total count
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	// Add ordering and pagination
	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Execute query
	var jobs []*model.Job
	err = r.db.SelectContext(ctx, &jobs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list jobs: %w", err)
	}

	return jobs, total, nil
}

// LockNext atomically locks the next available job for processing
func (r *JobRepository) LockNext(ctx context.Context, workerID string) (*model.Job, error) {
	var job model.Job

	query := `
		UPDATE jobs
		SET 
			status = $1,
			locked_at = $2,
			locked_by = $3,
			started_at = COALESCE(started_at, $2),
			attempts = attempts + 1
		WHERE id = (
			SELECT id
			FROM jobs
			WHERE status = $4
			  AND (next_retry_at IS NULL OR next_retry_at <= $2)
			ORDER BY created_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING 
			id, dataset_id, type, status, payload,
			attempts, max_attempts, last_error, next_retry_at,
			locked_at, locked_by, created_at, started_at, completed_at
	`

	now := time.Now().UTC()
	err := r.db.GetContext(ctx, &job, query,
		model.JobStatusRunning,
		now,
		workerID,
		model.JobStatusQueued,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No jobs available
		}
		return nil, fmt.Errorf("failed to lock job: %w", err)
	}

	return &job, nil
}

// UpdateStatus updates the status of a job
func (r *JobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.JobStatus) error {
	query := `
		UPDATE jobs
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

// GetDatasetCountryCode retrieves the country code for a job's dataset
func (r *JobRepository) GetDatasetCountryCode(ctx context.Context, jobID uuid.UUID) (string, error) {
	var countryCode string

	query := `
		SELECT c.code
		FROM jobs j
		JOIN datasets d ON d.id = j.dataset_id
		JOIN countries c ON c.id = d.country_id
		WHERE j.id = $1
	`

	err := r.db.GetContext(ctx, &countryCode, query, jobID)
	if err != nil {
		return "", fmt.Errorf("failed to get country code: %w", err)
	}

	return countryCode, nil
}

// Complete marks a job as completed
func (r *JobRepository) Complete(ctx context.Context, id uuid.UUID, result model.JobPayload) error {
	now := time.Now().UTC()

	query := `
		UPDATE jobs
		SET 
			status = $1,
			completed_at = $2,
			payload = $3,
			last_error = NULL
		WHERE id = $4
	`

	dbResult, err := r.db.ExecContext(ctx, query, model.JobStatusCompleted, now, result, id)
	if err != nil {
		return fmt.Errorf("failed to complete job: %w", err)
	}

	rows, err := dbResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

// Fail marks a job as failed
func (r *JobRepository) Fail(ctx context.Context, id uuid.UUID, errorMsg string) error {
	now := time.Now().UTC()

	query := `
		UPDATE jobs
		SET 
			status = $1,
			last_error = $2,
			completed_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, model.JobStatusFailed, errorMsg, now, id)
	if err != nil {
		return fmt.Errorf("failed to fail job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

// Retry marks a job for retry
func (r *JobRepository) Retry(ctx context.Context, id uuid.UUID, errorMsg string, nextRetry time.Time) error {
	query := `
		UPDATE jobs
		SET 
			status = $1,
			last_error = $2,
			next_retry_at = $3,
			locked_at = NULL,
			locked_by = NULL
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, model.JobStatusQueued, errorMsg, nextRetry, id)
	if err != nil {
		return fmt.Errorf("failed to retry job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}
