package parser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
)

// BatchInserter handles batch insertion of parsed records into the database
type BatchInserter struct {
	nameRepo  *repository.NameRepository
	batchSize int
}

// NewBatchInserter creates a new batch inserter
func NewBatchInserter(nameRepo *repository.NameRepository) *BatchInserter {
	return &BatchInserter{
		nameRepo:  nameRepo,
		batchSize: 1000, // Default batch size
	}
}

// SetBatchSize sets the batch size for insertions
func (b *BatchInserter) SetBatchSize(size int) {
	if size > 0 {
		b.batchSize = size
	}
}

// Insert processes records from a channel and inserts them in batches
// Returns the total number of records inserted
func (b *BatchInserter) Insert(
	ctx context.Context,
	datasetID uuid.UUID,
	countryID uuid.UUID,
	year int,
	records <-chan Record,
	errors <-chan error,
) (int, error) {
	totalInserted := 0
	batch := make([]*model.NameRecord, 0, b.batchSize)

	// Process records from channel
	for {
		select {
		case <-ctx.Done():
			return totalInserted, ctx.Err()

		case err, ok := <-errors:
			if ok && err != nil {
				return totalInserted, fmt.Errorf("parsing error: %w", err)
			}

		case record, ok := <-records:
			if !ok {
				// Channel closed, insert remaining batch
				if len(batch) > 0 {
					count, err := b.insertBatch(ctx, datasetID, countryID, year, batch)
					if err != nil {
						return totalInserted, err
					}
					totalInserted += count
				}
				return totalInserted, nil
			}

			// Add record to batch
			batch = append(batch, &model.NameRecord{
				Year:   record.Year,
				Name:   record.Name,
				Gender: record.Gender,
				Count:  record.Count,
			})

			// Insert batch when it reaches the batch size
			if len(batch) >= b.batchSize {
				count, err := b.insertBatch(ctx, datasetID, countryID, year, batch)
				if err != nil {
					return totalInserted, err
				}
				totalInserted += count
				batch = make([]*model.NameRecord, 0, b.batchSize)
			}
		}
	}
}

// insertBatch inserts a batch of records using the repository
func (b *BatchInserter) insertBatch(
	ctx context.Context,
	datasetID uuid.UUID,
	countryID uuid.UUID,
	year int,
	batch []*model.NameRecord,
) (int, error) {
	if len(batch) == 0 {
		return 0, nil
	}

	count, err := b.nameRepo.BatchInsert(ctx, datasetID, countryID, year, batch)
	if err != nil {
		return 0, fmt.Errorf("failed to insert batch: %w", err)
	}

	return count, nil
}

// InsertAll is a convenience method that collects all records and inserts them
// This is useful for smaller datasets or testing
func (b *BatchInserter) InsertAll(
	ctx context.Context,
	datasetID uuid.UUID,
	countryID uuid.UUID,
	year int,
	records []Record,
) (int, error) {
	if len(records) == 0 {
		return 0, nil
	}

	// Convert to NameRecords
	nameRecords := make([]*model.NameRecord, len(records))
	for i, record := range records {
		nameRecords[i] = &model.NameRecord{
			Year:   record.Year,
			Name:   record.Name,
			Gender: record.Gender,
			Count:  record.Count,
		}
	}

	// Insert in batches
	totalInserted := 0
	for i := 0; i < len(nameRecords); i += b.batchSize {
		end := i + b.batchSize
		if end > len(nameRecords) {
			end = len(nameRecords)
		}

		batch := nameRecords[i:end]
		count, err := b.nameRepo.BatchInsert(ctx, datasetID, countryID, year, batch)
		if err != nil {
			return totalInserted, fmt.Errorf("failed to insert batch: %w", err)
		}
		totalInserted += count
	}

	return totalInserted, nil
}
