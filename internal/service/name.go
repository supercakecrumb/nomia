package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
)

var (
	// ErrNameNotFound is returned when a name is not found
	ErrNameNotFound = errors.New("name not found")
	// ErrInvalidSearchQuery is returned when search query is invalid
	ErrInvalidSearchQuery = errors.New("search query must be at least 2 characters")
)

// NameService handles business logic for name operations
type NameService struct {
	nameRepo    *repository.NameRepository
	countryRepo *repository.CountryRepository
}

// NewNameService creates a new name service
func NewNameService(nameRepo *repository.NameRepository, countryRepo *repository.CountryRepository) *NameService {
	return &NameService{
		nameRepo:    nameRepo,
		countryRepo: countryRepo,
	}
}

// List retrieves names with filters, pagination, and rank calculation
func (s *NameService) List(
	ctx context.Context,
	countryCode string,
	year int,
	gender *string,
	namePrefix *string,
	minCount *int,
	sortBy string,
	sortOrder string,
	limit, offset int,
) ([]*model.NameResponse, int, error) {
	// Validate required parameters
	if countryCode == "" {
		return nil, 0, fmt.Errorf("%w: country is required", ErrInvalidInput)
	}
	if year < 1970 || year > 2030 {
		return nil, 0, fmt.Errorf("%w: year must be between 1970 and 2030", ErrInvalidInput)
	}

	// Validate optional parameters
	if gender != nil && *gender != "M" && *gender != "F" {
		return nil, 0, fmt.Errorf("%w: gender must be M or F", ErrInvalidInput)
	}
	if namePrefix != nil && len(*namePrefix) < 2 {
		return nil, 0, fmt.Errorf("%w: name prefix must be at least 2 characters", ErrInvalidInput)
	}
	if limit < 1 || limit > 1000 {
		return nil, 0, fmt.Errorf("%w: limit must be between 1 and 1000", ErrInvalidInput)
	}
	if offset < 0 {
		return nil, 0, fmt.Errorf("%w: offset must be >= 0", ErrInvalidInput)
	}

	// Get country by code
	country, err := s.countryRepo.GetByCode(ctx, countryCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrCountryNotFound
		}
		return nil, 0, fmt.Errorf("failed to get country: %w", err)
	}

	// Build filters
	filters := &model.NameFilters{
		CountryID: &country.ID,
		Year:      &year,
		Gender:    gender,
		Name:      namePrefix,
		MinCount:  minCount,
	}

	// Get names from repository
	names, total, err := s.nameRepo.ListWithRank(ctx, filters, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list names: %w", err)
	}

	// Calculate ranks for the filtered set
	ranks, err := s.nameRepo.CalculateRanks(ctx, country.ID, year, gender)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to calculate ranks: %w", err)
	}

	// Convert to response format
	responses := make([]*model.NameResponse, len(names))
	for i, name := range names {
		rank := ranks[name.Name]
		responses[i] = &model.NameResponse{
			Name:        name.Name,
			Gender:      name.Gender,
			Count:       name.Count,
			Year:        name.Year,
			CountryCode: country.Code,
			Rank:        rank,
		}
	}

	return responses, total, nil
}

// Search performs prefix search across years and countries
func (s *NameService) Search(
	ctx context.Context,
	query string,
	countryCode *string,
	gender *string,
	limit, offset int,
) ([]*model.NameSearchResponse, int, error) {
	// Validate query
	if len(query) < 2 {
		return nil, 0, ErrInvalidSearchQuery
	}

	// Validate optional parameters
	if gender != nil && *gender != "M" && *gender != "F" {
		return nil, 0, fmt.Errorf("%w: gender must be M or F", ErrInvalidInput)
	}
	if limit < 1 || limit > 1000 {
		return nil, 0, fmt.Errorf("%w: limit must be between 1 and 1000", ErrInvalidInput)
	}
	if offset < 0 {
		return nil, 0, fmt.Errorf("%w: offset must be >= 0", ErrInvalidInput)
	}

	// Get country ID if code provided
	var countryID *uuid.UUID
	if countryCode != nil {
		country, err := s.countryRepo.GetByCode(ctx, *countryCode)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, 0, ErrCountryNotFound
			}
			return nil, 0, fmt.Errorf("failed to get country: %w", err)
		}
		countryID = &country.ID
	}

	// Search names
	results, total, err := s.nameRepo.Search(ctx, query, countryID, gender, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search names: %w", err)
	}

	// Convert to response format
	responses := make([]*model.NameSearchResponse, len(results))
	for i, result := range results {
		responses[i] = &model.NameSearchResponse{
			Name:       result.Name,
			TotalCount: result.TotalCount,
			Countries:  result.Countries,
			YearRange: model.YearRange{
				Min: result.MinYear,
				Max: result.MaxYear,
			},
			GenderDist: model.GenderDistribution{
				M: result.MalePercentage,
				F: result.FemalePercentage,
			},
		}
	}

	return responses, total, nil
}

// GetByName retrieves detailed information about a specific name
func (s *NameService) GetByName(
	ctx context.Context,
	name string,
	countryCode *string,
) (*model.NameDetailResponse, error) {
	// Build filters
	var filters *model.NameFilters
	var country *model.Country
	var err error

	if countryCode != nil {
		country, err = s.countryRepo.GetByCode(ctx, *countryCode)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCountryNotFound
			}
			return nil, fmt.Errorf("failed to get country: %w", err)
		}
		filters = &model.NameFilters{
			CountryID: &country.ID,
		}
	}

	// Get all records for the name
	names, err := s.nameRepo.GetByName(ctx, name, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get name records: %w", err)
	}

	if len(names) == 0 {
		return nil, ErrNameNotFound
	}

	// Aggregate data
	var totalCount int64
	maleCount := int64(0)
	femaleCount := int64(0)
	peakYear := 0
	peakCount := 0
	countryMap := make(map[uuid.UUID]*model.CountryDetail)
	minYear := 9999
	maxYear := 0

	for _, n := range names {
		totalCount += int64(n.Count)

		if n.Gender == "M" {
			maleCount += int64(n.Count)
		} else {
			femaleCount += int64(n.Count)
		}

		if n.Count > peakCount {
			peakCount = n.Count
			peakYear = n.Year
		}

		if n.Year < minYear {
			minYear = n.Year
		}
		if n.Year > maxYear {
			maxYear = n.Year
		}

		// Aggregate by country
		if _, exists := countryMap[n.CountryID]; !exists {
			// Get country details
			c, err := s.countryRepo.GetByID(ctx, n.CountryID)
			if err == nil {
				countryMap[n.CountryID] = &model.CountryDetail{
					Code:  c.Code,
					Name:  c.Name,
					Count: 0,
					YearRange: model.YearRange{
						Min: n.Year,
						Max: n.Year,
					},
				}
			}
		}
		if cd, exists := countryMap[n.CountryID]; exists {
			cd.Count += int64(n.Count)
			if n.Year < cd.YearRange.Min {
				cd.YearRange.Min = n.Year
			}
			if n.Year > cd.YearRange.Max {
				cd.YearRange.Max = n.Year
			}
		}
	}

	// Convert country map to slice
	countries := make([]model.CountryDetail, 0, len(countryMap))
	for _, cd := range countryMap {
		countries = append(countries, *cd)
	}

	// Calculate gender distribution
	malePercentage := 0.0
	femalePercentage := 0.0
	if totalCount > 0 {
		malePercentage = float64(maleCount) / float64(totalCount) * 100
		femalePercentage = float64(femaleCount) / float64(totalCount) * 100
	}

	// Determine popularity trend (simplified)
	trend := "stable"
	if len(names) >= 2 {
		// Compare first half vs second half
		midpoint := len(names) / 2
		firstHalfSum := 0
		secondHalfSum := 0
		for i := 0; i < midpoint; i++ {
			firstHalfSum += names[i].Count
		}
		for i := midpoint; i < len(names); i++ {
			secondHalfSum += names[i].Count
		}
		if float64(secondHalfSum) > float64(firstHalfSum)*1.2 {
			trend = "increasing"
		} else if float64(secondHalfSum) < float64(firstHalfSum)*0.8 {
			trend = "decreasing"
		}
	}

	response := &model.NameDetailResponse{
		Name:       name,
		TotalCount: totalCount,
		Countries:  countries,
		GenderDistribution: model.GenderDistributionDetail{
			M: model.GenderStats{
				Count:      maleCount,
				Percentage: malePercentage,
			},
			F: model.GenderStats{
				Count:      femaleCount,
				Percentage: femalePercentage,
			},
		},
		PopularityTrend: trend,
		PeakYear:        peakYear,
		PeakCount:       peakCount,
	}

	return response, nil
}

// GetStats returns statistics for a dataset
func (s *NameService) GetStats(ctx context.Context, datasetID uuid.UUID) (map[string]interface{}, error) {
	return s.nameRepo.GetStats(ctx, datasetID)
}
