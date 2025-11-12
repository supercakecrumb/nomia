package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
)

var (
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
	// ErrCountryNotFound is returned when a country is not found
	ErrCountryNotFound = repository.ErrCountryNotFound
	// ErrCountryCodeExists is returned when a country code already exists
	ErrCountryCodeExists = repository.ErrCountryCodeExists
)

// CountryService handles business logic for countries
type CountryService struct {
	repo      *repository.CountryRepository
	validator *validator.Validate
}

// NewCountryService creates a new country service
func NewCountryService(repo *repository.CountryRepository) *CountryService {
	return &CountryService{
		repo:      repo,
		validator: validator.New(),
	}
}

// Create creates a new country
func (s *CountryService) Create(ctx context.Context, req *model.CreateCountryRequest) (*model.Country, error) {
	// Validate input
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Normalize code to uppercase
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	req.Name = strings.TrimSpace(req.Name)

	// Create country model
	country := &model.Country{
		ID:          uuid.New(),
		Code:        req.Code,
		Name:        req.Name,
		SourceURL:   req.SourceURL,
		Attribution: req.Attribution,
	}

	// Create in repository
	if err := s.repo.Create(ctx, country); err != nil {
		if errors.Is(err, repository.ErrCountryCodeExists) {
			return nil, ErrCountryCodeExists
		}
		return nil, fmt.Errorf("failed to create country: %w", err)
	}

	return country, nil
}

// GetByID retrieves a country by its ID
func (s *CountryService) GetByID(ctx context.Context, id uuid.UUID) (*model.Country, error) {
	country, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return country, nil
}

// GetByCode retrieves a country by its code
func (s *CountryService) GetByCode(ctx context.Context, code string) (*model.Country, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	country, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return country, nil
}

// GetByIDOrCode retrieves a country by its ID or code
func (s *CountryService) GetByIDOrCode(ctx context.Context, idOrCode string) (*model.Country, error) {
	// Try to parse as UUID first
	if id, err := uuid.Parse(idOrCode); err == nil {
		return s.GetByID(ctx, id)
	}

	// Otherwise, treat as code
	return s.GetByCode(ctx, idOrCode)
}

// GetWithStats retrieves a country with its statistics
func (s *CountryService) GetWithStats(ctx context.Context, idOrCode string) (*model.CountryWithStats, error) {
	// Get country
	country, err := s.GetByIDOrCode(ctx, idOrCode)
	if err != nil {
		return nil, err
	}

	// Get stats
	stats, err := s.repo.GetStats(ctx, country.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}

	return &model.CountryWithStats{
		Country: *country,
		Stats:   stats,
	}, nil
}

// List retrieves a paginated list of countries
func (s *CountryService) List(ctx context.Context, limit, offset int) ([]*model.Country, int, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	countries, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list countries: %w", err)
	}

	return countries, total, nil
}

// Update updates an existing country
func (s *CountryService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateCountryRequest) (*model.Country, error) {
	// Validate input
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Get existing country
	country, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		country.Name = strings.TrimSpace(*req.Name)
	}
	if req.SourceURL != nil {
		country.SourceURL = req.SourceURL
	}
	if req.Attribution != nil {
		country.Attribution = req.Attribution
	}

	// Update in repository
	if err := s.repo.Update(ctx, country); err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to update country: %w", err)
	}

	return country, nil
}

// Delete deletes a country by its ID
func (s *CountryService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return ErrCountryNotFound
		}
		return fmt.Errorf("failed to delete country: %w", err)
	}

	return nil
}
