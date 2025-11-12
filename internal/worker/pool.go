package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Pool manages a pool of worker goroutines
type Pool struct {
	processor    *Processor
	concurrency  int
	pollInterval time.Duration
	logger       *slog.Logger
	wg           sync.WaitGroup
	stopChan     chan struct{}
	stopped      bool
	mu           sync.Mutex
}

// NewPool creates a new worker pool
func NewPool(processor *Processor, concurrency int, pollInterval time.Duration, logger *slog.Logger) *Pool {
	return &Pool{
		processor:    processor,
		concurrency:  concurrency,
		pollInterval: pollInterval,
		logger:       logger,
		stopChan:     make(chan struct{}),
	}
}

// Start starts the worker pool
func (p *Pool) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return fmt.Errorf("pool has been stopped and cannot be restarted")
	}
	p.mu.Unlock()

	p.logger.Info("Starting worker pool",
		"concurrency", p.concurrency,
		"poll_interval", p.pollInterval,
	)

	// Start worker goroutines
	for i := 1; i <= p.concurrency; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	p.logger.Info("Worker pool started", "workers", p.concurrency)
	return nil
}

// Stop gracefully stops the worker pool
func (p *Pool) Stop() {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	p.mu.Unlock()

	p.logger.Info("Stopping worker pool")
	close(p.stopChan)
	p.wg.Wait()
	p.logger.Info("Worker pool stopped")
}

// worker is the main worker goroutine
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	workerID := fmt.Sprintf("worker-%d", id)
	logger := p.logger.With("worker_id", workerID)

	logger.Info("Worker started")
	defer logger.Info("Worker stopped")

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker context cancelled")
			return
		case <-p.stopChan:
			logger.Info("Worker received stop signal")
			return
		case <-ticker.C:
			// Try to lock and process a job
			if err := p.processNextJob(ctx, workerID, logger); err != nil {
				// Log error but continue processing
				if err != ErrNoJobsAvailable {
					logger.Error("Failed to process job", "error", err)
				}
			}
		}
	}
}

// processNextJob attempts to lock and process the next available job
func (p *Pool) processNextJob(ctx context.Context, workerID string, logger *slog.Logger) error {
	// Lock next job
	job, err := p.processor.LockNextJob(ctx, workerID)
	if err != nil {
		return err
	}

	if job == nil {
		// No jobs available
		return ErrNoJobsAvailable
	}

	logger.Info("Locked job for processing",
		"job_id", job.ID,
		"job_type", job.Type,
		"dataset_id", job.DatasetID,
		"attempt", job.Attempts,
	)

	// Process the job
	if err := p.processor.ProcessJob(ctx, job); err != nil {
		logger.Error("Job processing failed",
			"job_id", job.ID,
			"error", err,
		)
		return err
	}

	logger.Info("Job processed successfully",
		"job_id", job.ID,
		"dataset_id", job.DatasetID,
	)

	return nil
}

// ErrNoJobsAvailable is returned when no jobs are available for processing
var ErrNoJobsAvailable = fmt.Errorf("no jobs available")
