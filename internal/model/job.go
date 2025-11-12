package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the status of a background job
type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// JobType represents the type of background job
type JobType string

const (
	JobTypeParseDataset     JobType = "parse_dataset"
	JobTypeReprocessDataset JobType = "reprocess_dataset"
)

// Job represents a background job
type Job struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	DatasetID   *uuid.UUID `json:"dataset_id,omitempty" db:"dataset_id"`
	Type        JobType    `json:"type" db:"type"`
	Status      JobStatus  `json:"status" db:"status"`
	Payload     JobPayload `json:"payload,omitempty" db:"payload"`
	Attempts    int        `json:"attempts" db:"attempts"`
	MaxAttempts int        `json:"max_attempts" db:"max_attempts"`
	LastError   *string    `json:"last_error,omitempty" db:"last_error"`
	NextRetryAt *time.Time `json:"next_retry_at,omitempty" db:"next_retry_at"`
	LockedAt    *time.Time `json:"locked_at,omitempty" db:"locked_at"`
	LockedBy    *string    `json:"locked_by,omitempty" db:"locked_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// JobPayload represents the payload data for a job
type JobPayload map[string]interface{}

// Value implements the driver.Valuer interface for JobPayload
func (p JobPayload) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for JobPayload
func (p *JobPayload) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, p)
}

// CreateJobRequest represents the request to create a new job
type CreateJobRequest struct {
	DatasetID   *uuid.UUID `json:"dataset_id,omitempty"`
	Type        JobType    `json:"type" validate:"required"`
	Payload     JobPayload `json:"payload,omitempty"`
	MaxAttempts int        `json:"max_attempts,omitempty"`
}

// JobFilters represents filters for listing jobs
type JobFilters struct {
	Status    *JobStatus `json:"status,omitempty"`
	Type      *JobType   `json:"type,omitempty"`
	DatasetID *uuid.UUID `json:"dataset_id,omitempty"`
}

// JobResult represents the result of a completed job
type JobResult struct {
	RowsProcessed int      `json:"rows_processed,omitempty"`
	RowsSkipped   int      `json:"rows_skipped,omitempty"`
	Errors        []string `json:"errors,omitempty"`
}

// IsValid checks if the job status is valid
func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusQueued, JobStatusRunning, JobStatusCompleted, JobStatusFailed:
		return true
	}
	return false
}

// IsValid checks if the job type is valid
func (t JobType) IsValid() bool {
	switch t {
	case JobTypeParseDataset, JobTypeReprocessDataset:
		return true
	}
	return false
}

// ProcessingTimeSeconds returns the processing time in seconds
func (j *Job) ProcessingTimeSeconds() *int {
	if j.CompletedAt == nil || j.StartedAt == nil {
		return nil
	}
	seconds := int(j.CompletedAt.Sub(*j.StartedAt).Seconds())
	return &seconds
}
