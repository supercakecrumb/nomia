package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface for local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// Save saves a file to local filesystem
// Creates directory structure: basePath/{datasetID}/original.csv
func (s *LocalStorage) Save(ctx context.Context, datasetID string, filename string, reader io.Reader) (string, error) {
	// Create dataset directory
	datasetDir := filepath.Join(s.basePath, datasetID)
	if err := os.MkdirAll(datasetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create dataset directory: %w", err)
	}

	// Use "original.csv" as the stored filename
	storedFilename := "original.csv"
	filePath := filepath.Join(datasetDir, storedFilename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data from reader to file
	if _, err := io.Copy(file, reader); err != nil {
		// Clean up on error
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative path from base
	relativePath := filepath.Join(datasetID, storedFilename)
	return relativePath, nil
}

// Load retrieves a file from local filesystem
func (s *LocalStorage) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from local filesystem
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted, not an error
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Try to remove the parent directory if it's empty
	parentDir := filepath.Dir(fullPath)
	os.Remove(parentDir) // Ignore error, directory might not be empty

	return nil
}

// Exists checks if a file exists in local filesystem
func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}
