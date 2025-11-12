package model

import (
	"time"

	"github.com/google/uuid"
)

// Name represents a baby name record in the database
type Name struct {
	ID        int64     `db:"id" json:"id"`
	DatasetID uuid.UUID `db:"dataset_id" json:"dataset_id"`
	CountryID uuid.UUID `db:"country_id" json:"country_id"`
	Year      int       `db:"year" json:"year"`
	Name      string    `db:"name" json:"name"`
	Gender    string    `db:"gender" json:"gender"` // M or F
	Count     int       `db:"count" json:"count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NameRecord represents a parsed name record before database insertion
// This is the intermediate format used during CSV parsing
type NameRecord struct {
	Year   int    `json:"year"`
	Name   string `json:"name"`
	Gender string `json:"gender"` // M or F
	Count  int    `json:"count"`
}

// NameFilters represents filters for querying names
type NameFilters struct {
	DatasetID *uuid.UUID `json:"dataset_id,omitempty"`
	CountryID *uuid.UUID `json:"country_id,omitempty"`
	Year      *int       `json:"year,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	MinCount  *int       `json:"min_count,omitempty"`
}

// NameResponse represents a name record in API responses
type NameResponse struct {
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Count       int    `json:"count"`
	Year        int    `json:"year"`
	CountryCode string `json:"country_code"`
	Rank        int    `json:"rank"`
}

// NameSearchResult represents aggregated search results across years
type NameSearchResult struct {
	Name             string   `json:"name" db:"name"`
	TotalCount       int64    `json:"total_count" db:"total_count"`
	Countries        []string `json:"countries" db:"countries"`
	MinYear          int      `json:"min_year" db:"min_year"`
	MaxYear          int      `json:"max_year" db:"max_year"`
	PrimaryGender    *string  `json:"primary_gender,omitempty" db:"primary_gender"`
	MalePercentage   float64  `json:"male_percentage" db:"male_percentage"`
	FemalePercentage float64  `json:"female_percentage" db:"female_percentage"`
}

// NameSearchResponse represents the API response for search
type NameSearchResponse struct {
	Name       string             `json:"name"`
	TotalCount int64              `json:"total_count"`
	Countries  []string           `json:"countries"`
	YearRange  YearRange          `json:"year_range"`
	GenderDist GenderDistribution `json:"gender_distribution"`
}

// YearRange represents a year range
type YearRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// GenderDistribution represents gender distribution percentages
type GenderDistribution struct {
	M float64 `json:"M"`
	F float64 `json:"F"`
}

// NameDetailResponse represents detailed information about a name
type NameDetailResponse struct {
	Name               string                   `json:"name"`
	TotalCount         int64                    `json:"total_count"`
	Countries          []CountryDetail          `json:"countries"`
	GenderDistribution GenderDistributionDetail `json:"gender_distribution"`
	PopularityTrend    string                   `json:"popularity_trend"`
	PeakYear           int                      `json:"peak_year"`
	PeakCount          int                      `json:"peak_count"`
}

// CountryDetail represents country-specific name details
type CountryDetail struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Count     int64     `json:"count"`
	YearRange YearRange `json:"year_range"`
}

// GenderDistributionDetail represents detailed gender distribution
type GenderDistributionDetail struct {
	M GenderStats `json:"M"`
	F GenderStats `json:"F"`
}

// GenderStats represents statistics for a gender
type GenderStats struct {
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}
