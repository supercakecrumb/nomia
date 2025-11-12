package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/database"
	"github.com/supercakecrumb/affirm-name/internal/logging"
	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/service"
	"github.com/supercakecrumb/affirm-name/internal/storage"
	"github.com/supercakecrumb/affirm-name/internal/worker"
)

func TestWorkerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Initialize logger
	logger := logging.InitLogger("debug", "text")

	// Connect to database
	db, err := database.New(cfg.Database)
	require.NoError(t, err)
	defer db.Close()

	// Initialize storage
	tempDir := t.TempDir()
	stor, err := storage.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Initialize repositories
	countryRepo := repository.NewCountryRepository(db.Pool)
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	jobRepo := repository.NewJobRepository(db.Sqlx)
	nameRepo := repository.NewNameRepository(db.Sqlx)

	// Initialize services
	parserService := service.NewParserService(datasetRepo, nameRepo, jobRepo, stor, logger)

	// Create test country
	ctx := context.Background()
	country := &model.Country{
		ID:   uuid.New(),
		Code: "US",
		Name: "United States",
	}
	err = countryRepo.Create(ctx, country)
	require.NoError(t, err)
	defer countryRepo.Delete(ctx, country.ID)

	// Copy test file to storage
	testFile := filepath.Join("fixtures", "us_ssa_sample.csv")
	testData, err := os.ReadFile(testFile)
	require.NoError(t, err)

	storagePath := "test_" + uuid.New().String() + ".csv"
	err = os.WriteFile(filepath.Join(tempDir, storagePath), testData, 0644)
	require.NoError(t, err)

	// Create test dataset
	dataset, err := datasetRepo.Create(ctx, &model.CreateDatasetRequest{
		CountryID:  country.ID,
		Filename:   "yob2023.txt",
		FilePath:   storagePath,
		FileSize:   int64(len(testData)),
		UploadedBy: "test-user",
	})
	require.NoError(t, err)

	// Create test job
	job, err := jobRepo.Create(ctx, &model.CreateJobRequest{
		DatasetID: &dataset.ID,
		Type:      model.JobTypeParseDataset,
		Payload: model.JobPayload{
			"dataset_id": dataset.ID.String(),
		},
		MaxAttempts: 3,
	})
	require.NoError(t, err)

	// Initialize processor
	processor := worker.NewProcessor(
		jobRepo,
		datasetRepo,
		parserService,
		stor,
		3,
		logger,
	)

	// Test: Lock next job
	t.Run("LockNextJob", func(t *testing.T) {
		lockedJob, err := processor.LockNextJob(ctx, "test-worker")
		require.NoError(t, err)
		require.NotNil(t, lockedJob)
		assert.Equal(t, job.ID, lockedJob.ID)
		assert.Equal(t, model.JobStatusRunning, lockedJob.Status)
		assert.Equal(t, 1, lockedJob.Attempts)
	})

	// Test: Process job
	t.Run("ProcessJob", func(t *testing.T) {
		// Get the locked job
		lockedJob, err := jobRepo.GetByID(ctx, job.ID)
		require.NoError(t, err)

		// Process the job
		err = processor.ProcessJob(ctx, lockedJob)
		require.NoError(t, err)

		// Verify job status
		updatedJob, err := jobRepo.GetByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, model.JobStatusCompleted, updatedJob.Status)
		assert.NotNil(t, updatedJob.CompletedAt)

		// Verify dataset status
		updatedDataset, err := datasetRepo.GetByID(ctx, dataset.ID)
		require.NoError(t, err)
		assert.Equal(t, model.DatasetStatusCompleted, updatedDataset.Status)
		assert.NotNil(t, updatedDataset.RowCount)
		assert.Greater(t, *updatedDataset.RowCount, 0)
	})

	// Test: Worker pool
	t.Run("WorkerPool", func(t *testing.T) {
		// Create another job
		dataset2, err := datasetRepo.Create(ctx, &model.CreateDatasetRequest{
			CountryID:  country.ID,
			Filename:   "yob2024.txt",
			FilePath:   storagePath,
			FileSize:   int64(len(testData)),
			UploadedBy: "test-user",
		})
		require.NoError(t, err)

		job2, err := jobRepo.Create(ctx, &model.CreateJobRequest{
			DatasetID: &dataset2.ID,
			Type:      model.JobTypeParseDataset,
			Payload: model.JobPayload{
				"dataset_id": dataset2.ID.String(),
			},
			MaxAttempts: 3,
		})
		require.NoError(t, err)

		// Create worker pool
		pool := worker.NewPool(
			processor,
			2,                    // 2 workers
			100*time.Millisecond, // Fast polling for test
			logger,
		)

		// Start pool
		poolCtx, cancel := context.WithCancel(ctx)
		err = pool.Start(poolCtx)
		require.NoError(t, err)

		// Wait for job to be processed
		time.Sleep(500 * time.Millisecond)

		// Stop pool
		cancel()
		pool.Stop()

		// Verify job was processed
		processedJob, err := jobRepo.GetByID(ctx, job2.ID)
		require.NoError(t, err)
		assert.Equal(t, model.JobStatusCompleted, processedJob.Status)
	})
}

func TestWorkerRetryLogic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Initialize logger
	logger := logging.InitLogger("debug", "text")

	// Connect to database
	db, err := database.New(cfg.Database)
	require.NoError(t, err)
	defer db.Close()

	// Initialize storage
	tempDir := t.TempDir()
	stor, err := storage.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Initialize repositories
	countryRepo := repository.NewCountryRepository(db.Pool)
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	jobRepo := repository.NewJobRepository(db.Sqlx)
	nameRepo := repository.NewNameRepository(db.Sqlx)

	// Initialize services
	parserService := service.NewParserService(datasetRepo, nameRepo, jobRepo, stor, logger)

	// Create test country
	ctx := context.Background()
	country := &model.Country{
		ID:   uuid.New(),
		Code: "US",
		Name: "United States",
	}
	err = countryRepo.Create(ctx, country)
	require.NoError(t, err)
	defer countryRepo.Delete(ctx, country.ID)

	// Create test dataset with invalid file path (will cause error)
	dataset, err := datasetRepo.Create(ctx, &model.CreateDatasetRequest{
		CountryID:  country.ID,
		Filename:   "yob2023.txt",
		FilePath:   "nonexistent.csv",
		FileSize:   100,
		UploadedBy: "test-user",
	})
	require.NoError(t, err)

	// Create test job
	job, err := jobRepo.Create(ctx, &model.CreateJobRequest{
		DatasetID: &dataset.ID,
		Type:      model.JobTypeParseDataset,
		Payload: model.JobPayload{
			"dataset_id": dataset.ID.String(),
		},
		MaxAttempts: 3,
	})
	require.NoError(t, err)

	// Initialize processor
	processor := worker.NewProcessor(
		jobRepo,
		datasetRepo,
		parserService,
		stor,
		3,
		logger,
	)

	// Test: Job fails and is retried
	t.Run("JobRetry", func(t *testing.T) {
		// Lock and process job (will fail)
		lockedJob, err := processor.LockNextJob(ctx, "test-worker")
		require.NoError(t, err)
		require.NotNil(t, lockedJob)

		err = processor.ProcessJob(ctx, lockedJob)
		require.NoError(t, err) // No error because retry was scheduled

		// Verify job was scheduled for retry
		retriedJob, err := jobRepo.GetByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, model.JobStatusQueued, retriedJob.Status)
		assert.NotNil(t, retriedJob.NextRetryAt)
		assert.NotNil(t, retriedJob.LastError)
		assert.Equal(t, 1, retriedJob.Attempts)
	})
}
