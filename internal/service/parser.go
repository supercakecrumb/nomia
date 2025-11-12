package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/parser"
	_ "github.com/supercakecrumb/affirm-name/internal/parser/parsers" // Import parsers to register them
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/storage"
)

// ParserService handles parsing of dataset files
type ParserService struct {
	datasetRepo *repository.DatasetRepository
	nameRepo    *repository.NameRepository
	jobRepo     *repository.JobRepository
	storage     storage.Storage
	logger      *slog.Logger
}

// NewParserService creates a new parser service
func NewParserService(
	datasetRepo *repository.DatasetRepository,
	nameRepo *repository.NameRepository,
	jobRepo *repository.JobRepository,
	storage storage.Storage,
	logger *slog.Logger,
) *ParserService {
	return &ParserService{
		datasetRepo: datasetRepo,
		nameRepo:    nameRepo,
		jobRepo:     jobRepo,
		storage:     storage,
		logger:      logger,
	}
}

// ProcessDataset processes a dataset file by parsing and inserting records
func (s *ParserService) ProcessDataset(ctx context.Context, datasetID uuid.UUID) error {
	s.logger.Info("Starting dataset processing", "dataset_id", datasetID)

	// Get dataset with country information
	dataset, err := s.datasetRepo.GetByIDWithCountry(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to get dataset: %w", err)
	}

	// Update dataset status to processing
	if err := s.datasetRepo.UpdateStatus(ctx, datasetID, model.DatasetStatusProcessing); err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	// Get parser for country
	countryParser, err := parser.Get(dataset.CountryCode)
	if err != nil {
		s.logger.Error("No parser available for country", "country_code", dataset.CountryCode, "error", err)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("No parser available for country: %s", dataset.CountryCode)); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("no parser for country %s: %w", dataset.CountryCode, err)
	}

	// Extract year from filename (e.g., "yob2023.txt" -> 2023)
	year, err := s.extractYearFromFilename(dataset.Filename)
	if err != nil {
		s.logger.Error("Failed to extract year from filename", "filename", dataset.Filename, "error", err)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("Invalid filename format: %s", err.Error())); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("failed to extract year: %w", err)
	}

	s.logger.Info("Extracted year from filename", "year", year, "filename", dataset.Filename)

	// Open file from storage
	reader, err := s.storage.Load(ctx, dataset.FilePath)
	if err != nil {
		s.logger.Error("Failed to open file from storage", "file_path", dataset.FilePath, "error", err)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("Failed to open file: %s", err.Error())); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer reader.Close()

	// Validate file format
	if err := countryParser.Validate(reader); err != nil {
		s.logger.Error("File validation failed", "error", err)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("Invalid file format: %s", err.Error())); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("file validation failed: %w", err)
	}

	// Re-open file for parsing (validation consumed the reader)
	reader.Close()
	reader, err = s.storage.Load(ctx, dataset.FilePath)
	if err != nil {
		s.logger.Error("Failed to re-open file from storage", "error", err)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("Failed to re-open file: %s", err.Error())); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("failed to re-open file: %w", err)
	}
	defer reader.Close()

	// Parse CSV file
	records, errors := countryParser.Parse(ctx, reader)

	// Create batch inserter
	inserter := parser.NewBatchInserter(s.nameRepo)

	// Insert records
	s.logger.Info("Starting batch insert", "dataset_id", datasetID, "country_id", dataset.CountryID, "year", year)

	rowCount, err := inserter.Insert(ctx, datasetID, dataset.CountryID, year, records, errors)
	if err != nil {
		s.logger.Error("Failed to insert records", "error", err, "rows_inserted", rowCount)
		if updateErr := s.datasetRepo.UpdateFailed(ctx, datasetID, fmt.Sprintf("Failed to insert records: %s", err.Error())); updateErr != nil {
			s.logger.Error("Failed to update dataset status", "error", updateErr)
		}
		return fmt.Errorf("failed to insert records: %w", err)
	}

	s.logger.Info("Successfully inserted records", "row_count", rowCount)

	// Update dataset as completed
	if err := s.datasetRepo.UpdateCompleted(ctx, datasetID, rowCount); err != nil {
		s.logger.Error("Failed to update dataset as completed", "error", err)
		return fmt.Errorf("failed to update dataset: %w", err)
	}

	s.logger.Info("Dataset processing completed", "dataset_id", datasetID, "row_count", rowCount)
	return nil
}

// extractYearFromFilename extracts the year from a filename
// Expected formats: yob2023.txt, names_2023.csv, 2023.csv, etc.
func (s *ParserService) extractYearFromFilename(filename string) (int, error) {
	// Remove extension
	name := filename
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		name = filename[:idx]
	}

	// Try to find a 4-digit year in the filename
	for i := 0; i <= len(name)-4; i++ {
		substr := name[i : i+4]
		if year, err := strconv.Atoi(substr); err == nil {
			// Check if it's a reasonable year (1880-2100)
			if year >= 1880 && year <= 2100 {
				return year, nil
			}
		}
	}

	return 0, fmt.Errorf("no valid year found in filename: %s", filename)
}

// ReprocessDataset reprocesses a dataset by deleting existing records and re-parsing
func (s *ParserService) ReprocessDataset(ctx context.Context, datasetID uuid.UUID) error {
	s.logger.Info("Starting dataset reprocessing", "dataset_id", datasetID)

	// Update status to reprocessing
	if err := s.datasetRepo.UpdateStatus(ctx, datasetID, model.DatasetStatusReprocessing); err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	// Delete existing records
	s.logger.Info("Deleting existing records", "dataset_id", datasetID)
	if err := s.nameRepo.Delete(ctx, datasetID); err != nil {
		s.logger.Warn("No existing records to delete or delete failed", "error", err)
		// Continue anyway - might be first processing
	}

	// Process the dataset
	if err := s.ProcessDataset(ctx, datasetID); err != nil {
		return fmt.Errorf("failed to reprocess dataset: %w", err)
	}

	s.logger.Info("Dataset reprocessing completed", "dataset_id", datasetID)
	return nil
}

// ValidateFile validates a file without processing it
func (s *ParserService) ValidateFile(ctx context.Context, countryCode string, reader io.Reader) error {
	// Get parser for country
	countryParser, err := parser.Get(countryCode)
	if err != nil {
		return fmt.Errorf("no parser for country %s: %w", countryCode, err)
	}

	// Validate file format
	if err := countryParser.Validate(reader); err != nil {
		return fmt.Errorf("file validation failed: %w", err)
	}

	return nil
}

// GetAvailableParsers returns metadata for all available parsers
func (s *ParserService) GetAvailableParsers() []parser.ParserMetadata {
	return parser.GetMetadata()
}
