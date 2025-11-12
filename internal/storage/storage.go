package storage

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Save saves a file to storage and returns the storage path
	Save(ctx context.Context, datasetID string, filename string, reader io.Reader) (string, error)

	// Load retrieves a file from storage
	Load(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, path string) (bool, error)
}
