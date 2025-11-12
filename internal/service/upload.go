package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/storage"
)

const (
	// MaxFileSize is the maximum allowed file size (100MB)
	MaxFileSize = 100 * 1024 * 1024

	// AllowedMIMEType is the allowed MIME type for uploads
	AllowedMIMEType = "text/csv"
)

// UploadService handles file upload operations
type UploadService struct {
	storage     storage.Storage
	datasetRepo *repository.DatasetRepository
	jobRepo     *repository.JobRepository
	countryRepo *repository.CountryRepository
}

// NewUploadService creates a new UploadService
func NewUploadService(
	storage storage.Storage,
	datasetRepo *repository.DatasetRepository,
	jobRepo *repository.JobRepository,
	countryRepo *repository.CountryRepository,
) *UploadService {
	return &UploadService{
		storage:     storage,
		datasetRepo: datasetRepo,
		jobRepo:     jobRepo,
		countryRepo: countryRepo,
	}
}

// UploadRequest represents a file upload request
type UploadRequest struct {
	File       *multipart.FileHeader
	CountryID  uuid.UUID
	UploadedBy string
}

// UploadResponse represents the response from an upload operation
type UploadResponse struct {
	DatasetID uuid.UUID `json:"dataset_id"`
	JobID     uuid.UUID `json:"job_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

// Upload handles the file upload process
func (s *UploadService) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	// Validate file size
	if req.File.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	if req.File.Size == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	// Validate file type
	if err := s.validateFileType(req.File); err != nil {
		return nil, err
	}

	// Verify country exists
	_, err := s.countryRepo.GetByID(ctx, req.CountryID)
	if err != nil {
		return nil, fmt.Errorf("country not found: %w", err)
	}

	// Open the uploaded file
	file, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Generate dataset ID
	datasetID := uuid.New()

	// Save file to storage
	filePath, err := s.storage.Save(ctx, datasetID.String(), req.File.Filename, file)
	if err != nil {
		return nil, fmt.Errorf("failed to save file to storage: %w", err)
	}

	// Create dataset record
	dataset, err := s.datasetRepo.Create(ctx, &model.CreateDatasetRequest{
		CountryID:  req.CountryID,
		Filename:   req.File.Filename,
		FilePath:   filePath,
		FileSize:   req.File.Size,
		UploadedBy: req.UploadedBy,
	})
	if err != nil {
		// Try to clean up the uploaded file
		_ = s.storage.Delete(ctx, filePath)
		return nil, fmt.Errorf("failed to create dataset record: %w", err)
	}

	// Create job for processing
	job, err := s.jobRepo.Create(ctx, &model.CreateJobRequest{
		DatasetID: &dataset.ID,
		Type:      model.JobTypeParseDataset,
		Payload: model.JobPayload{
			"dataset_id": dataset.ID.String(),
			"country_id": req.CountryID.String(),
			"file_path":  filePath,
			"filename":   req.File.Filename,
		},
		MaxAttempts: 3,
	})
	if err != nil {
		// Dataset record exists but job creation failed
		// The dataset will remain in "pending" status
		return nil, fmt.Errorf("failed to create processing job: %w", err)
	}

	return &UploadResponse{
		DatasetID: dataset.ID,
		JobID:     job.ID,
		Status:    string(model.DatasetStatusPending),
		Message:   "Dataset uploaded successfully. Processing will begin shortly.",
	}, nil
}

// validateFileType validates the file type based on extension and MIME type
func (s *UploadService) validateFileType(file *multipart.FileHeader) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".csv" && ext != ".txt" {
		return fmt.Errorf("invalid file extension: %s (only .csv and .txt files are allowed)", ext)
	}

	// Check MIME type
	contentType := file.Header.Get("Content-Type")
	if contentType != "" && !isAllowedMIMEType(contentType) {
		return fmt.Errorf("invalid file type: %s (only CSV files are allowed)", contentType)
	}

	return nil
}

// isAllowedMIMEType checks if the MIME type is allowed
func isAllowedMIMEType(mimeType string) bool {
	allowedTypes := []string{
		"text/csv",
		"text/plain",
		"application/csv",
		"application/vnd.ms-excel",
	}

	for _, allowed := range allowedTypes {
		if strings.HasPrefix(mimeType, allowed) {
			return true
		}
	}

	return false
}

// ValidateFile performs pre-upload validation without saving the file
func (s *UploadService) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	return s.validateFileType(file)
}

// GetUploadLimits returns the upload limits
func (s *UploadService) GetUploadLimits() map[string]interface{} {
	return map[string]interface{}{
		"max_file_size":      MaxFileSize,
		"max_file_size_mb":   MaxFileSize / (1024 * 1024),
		"allowed_extensions": []string{".csv", ".txt"},
		"allowed_mime_types": []string{"text/csv", "text/plain", "application/csv"},
	}
}
