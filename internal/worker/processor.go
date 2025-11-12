package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/service"
	"github.com/supercakecrumb/affirm-name/internal/storage"
)

// Processor handles job processing logic
type Processor struct {
	jobRepo     *repository.JobRepository
	datasetRepo *repository.DatasetRepository
	parserSvc   *service.ParserService
	storage     storage.Storage
	maxRetries  int
	logger      *slog.Logger
}

// NewProcessor creates a new job processor
func NewProcessor(
	jobRepo *repository.JobRepository,
	datasetRepo *repository.DatasetRepository,
	parserSvc *service.ParserService,
	storage storage.Storage,
	maxRetries int,
	logger *slog.Logger,
) *Processor {
	return &Processor{
		jobRepo:     jobRepo,
		datasetRepo: datasetRepo,
		parserSvc:   parserSvc,
		storage:     storage,
		maxRetries:  maxRetries,
		logger:      logger,
	}
}

// LockNextJob atomically locks the next available job for processing
func (p *Processor) LockNextJob(ctx context.Context, workerID string) (*model.Job, error) {
	job, err := p.jobRepo.LockNext(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock job: %w", err)
	}

	if job == nil {
		return nil, nil // No jobs available
	}

	return job, nil
}

// ProcessJob processes a single job
func (p *Processor) ProcessJob(ctx context.Context, job *model.Job) error {
	logger := p.logger.With(
		"job_id", job.ID,
		"job_type", job.Type,
		"dataset_id", job.DatasetID,
		"attempt", job.Attempts,
	)

	logger.Info("Processing job")

	// Process based on job type
	var err error
	switch job.Type {
	case model.JobTypeParseDataset:
		err = p.processParseDataset(ctx, job, logger)
	case model.JobTypeReprocessDataset:
		err = p.processReprocessDataset(ctx, job, logger)
	default:
		err = fmt.Errorf("unknown job type: %s", job.Type)
	}

	if err != nil {
		logger.Error("Job processing failed", "error", err)
		return p.handleJobFailure(ctx, job, err, logger)
	}

	// Mark job as completed
	result := model.JobPayload{
		"completed_at": time.Now().UTC(),
	}

	if err := p.jobRepo.Complete(ctx, job.ID, result); err != nil {
		logger.Error("Failed to mark job as completed", "error", err)
		return fmt.Errorf("failed to complete job: %w", err)
	}

	logger.Info("Job completed successfully")
	return nil
}

// processParseDataset processes a parse_dataset job
func (p *Processor) processParseDataset(ctx context.Context, job *model.Job, logger *slog.Logger) error {
	if job.DatasetID == nil {
		return fmt.Errorf("dataset_id is required for parse_dataset job")
	}

	datasetID := *job.DatasetID
	logger.Info("Processing parse_dataset job", "dataset_id", datasetID)

	// Update dataset status to processing
	if err := p.datasetRepo.UpdateStatus(ctx, datasetID, model.DatasetStatusProcessing); err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	// Process the dataset using ParserService
	if err := p.parserSvc.ProcessDataset(ctx, datasetID); err != nil {
		// Update dataset status to failed
		if updateErr := p.datasetRepo.UpdateFailed(ctx, datasetID, err.Error()); updateErr != nil {
			logger.Error("Failed to update dataset status to failed", "error", updateErr)
		}
		return fmt.Errorf("failed to process dataset: %w", err)
	}

	logger.Info("Dataset processed successfully", "dataset_id", datasetID)
	return nil
}

// processReprocessDataset processes a reprocess_dataset job
func (p *Processor) processReprocessDataset(ctx context.Context, job *model.Job, logger *slog.Logger) error {
	if job.DatasetID == nil {
		return fmt.Errorf("dataset_id is required for reprocess_dataset job")
	}

	datasetID := *job.DatasetID
	logger.Info("Processing reprocess_dataset job", "dataset_id", datasetID)

	// Update dataset status to reprocessing
	if err := p.datasetRepo.UpdateStatus(ctx, datasetID, model.DatasetStatusReprocessing); err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	// Reprocess the dataset using ParserService
	if err := p.parserSvc.ReprocessDataset(ctx, datasetID); err != nil {
		// Update dataset status to failed
		if updateErr := p.datasetRepo.UpdateFailed(ctx, datasetID, err.Error()); updateErr != nil {
			logger.Error("Failed to update dataset status to failed", "error", updateErr)
		}
		return fmt.Errorf("failed to reprocess dataset: %w", err)
	}

	logger.Info("Dataset reprocessed successfully", "dataset_id", datasetID)
	return nil
}

// handleJobFailure handles job failure with retry logic
func (p *Processor) handleJobFailure(ctx context.Context, job *model.Job, jobErr error, logger *slog.Logger) error {
	errorMsg := jobErr.Error()

	// Check if we should retry
	if job.Attempts < p.maxRetries && p.shouldRetry(jobErr) {
		// Calculate next retry time with exponential backoff
		nextRetry := p.calculateNextRetry(job.Attempts)

		logger.Info("Scheduling job retry",
			"attempt", job.Attempts,
			"max_attempts", p.maxRetries,
			"next_retry_at", nextRetry,
		)

		// Schedule retry
		if err := p.jobRepo.Retry(ctx, job.ID, errorMsg, nextRetry); err != nil {
			logger.Error("Failed to schedule job retry", "error", err)
			return fmt.Errorf("failed to schedule retry: %w", err)
		}

		return nil
	}

	// Max retries reached or permanent error - fail the job
	logger.Error("Job failed permanently",
		"attempt", job.Attempts,
		"max_attempts", p.maxRetries,
		"error", errorMsg,
	)

	if err := p.jobRepo.Fail(ctx, job.ID, errorMsg); err != nil {
		logger.Error("Failed to mark job as failed", "error", err)
		return fmt.Errorf("failed to fail job: %w", err)
	}

	// Also update dataset status to failed if applicable
	if job.DatasetID != nil {
		if err := p.datasetRepo.UpdateFailed(ctx, *job.DatasetID, errorMsg); err != nil {
			logger.Error("Failed to update dataset status to failed", "error", err)
		}
	}

	return nil
}

// shouldRetry determines if an error is retryable
func (p *Processor) shouldRetry(err error) bool {
	// Check for permanent errors that should not be retried
	errStr := err.Error()

	// Permanent errors - do not retry
	permanentErrors := []string{
		"no parser",
		"invalid filename format",
		"invalid file format",
		"file validation failed",
		"dataset not found",
		"job not found",
	}

	for _, permErr := range permanentErrors {
		if contains(errStr, permErr) {
			return false
		}
	}

	// Transient errors - retry
	// Database connection errors, storage errors, etc.
	return true
}

// calculateNextRetry calculates the next retry time with exponential backoff
// Backoff schedule: 1min, 5min, 15min
func (p *Processor) calculateNextRetry(attempt int) time.Time {
	var delay time.Duration

	switch attempt {
	case 1:
		delay = 1 * time.Minute
	case 2:
		delay = 5 * time.Minute
	default:
		delay = 15 * time.Minute
	}

	return time.Now().UTC().Add(delay)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

// findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
