package model

import (
	"time"

	"github.com/google/uuid"
)

// DatasetStatus represents the processing status of a dataset
type DatasetStatus string

const (
	DatasetStatusPending      DatasetStatus = "pending"
	DatasetStatusProcessing   DatasetStatus = "processing"
	DatasetStatusCompleted    DatasetStatus = "completed"
	DatasetStatusFailed       DatasetStatus = "failed"
	DatasetStatusReprocessing DatasetStatus = "reprocessing"
)

// Dataset represents an uploaded dataset file
type Dataset struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	CountryID    uuid.UUID     `json:"country_id" db:"country_id"`
	Filename     string        `json:"filename" db:"filename"`
	FilePath     string        `json:"file_path" db:"file_path"`
	FileSize     int64         `json:"file_size" db:"file_size"`
	Status       DatasetStatus `json:"status" db:"status"`
	RowCount     *int          `json:"row_count,omitempty" db:"row_count"`
	ErrorMessage *string       `json:"error_message,omitempty" db:"error_message"`
	UploadedBy   string        `json:"uploaded_by" db:"uploaded_by"`
	UploadedAt   time.Time     `json:"uploaded_at" db:"uploaded_at"`
	ProcessedAt  *time.Time    `json:"processed_at,omitempty" db:"processed_at"`
	DeletedAt    *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CreateDatasetRequest represents the request to create a new dataset
type CreateDatasetRequest struct {
	CountryID  uuid.UUID `json:"country_id" validate:"required"`
	Filename   string    `json:"filename" validate:"required"`
	FilePath   string    `json:"file_path" validate:"required"`
	FileSize   int64     `json:"file_size" validate:"required,gt=0"`
	UploadedBy string    `json:"uploaded_by" validate:"required"`
}

// UpdateDatasetRequest represents the request to update a dataset
type UpdateDatasetRequest struct {
	Status       *DatasetStatus `json:"status,omitempty"`
	RowCount     *int           `json:"row_count,omitempty"`
	ErrorMessage *string        `json:"error_message,omitempty"`
	ProcessedAt  *time.Time     `json:"processed_at,omitempty"`
}

// DatasetFilters represents filters for listing datasets
type DatasetFilters struct {
	CountryID *uuid.UUID     `json:"country_id,omitempty"`
	Status    *DatasetStatus `json:"status,omitempty"`
}

// DatasetWithCountry represents a dataset with country information
type DatasetWithCountry struct {
	Dataset
	CountryCode string `json:"country_code" db:"country_code"`
	CountryName string `json:"country_name" db:"country_name"`
}

// IsValid checks if the dataset status is valid
func (s DatasetStatus) IsValid() bool {
	switch s {
	case DatasetStatusPending, DatasetStatusProcessing, DatasetStatusCompleted,
		DatasetStatusFailed, DatasetStatusReprocessing:
		return true
	}
	return false
}

// ProcessingTimeSeconds returns the processing time in seconds
func (d *Dataset) ProcessingTimeSeconds() *int {
	if d.ProcessedAt == nil {
		return nil
	}
	seconds := int(d.ProcessedAt.Sub(d.UploadedAt).Seconds())
	return &seconds
}
