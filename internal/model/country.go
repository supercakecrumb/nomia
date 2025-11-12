package model

import (
	"time"

	"github.com/google/uuid"
)

// Country represents a country with baby name statistics
type Country struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Code        string    `json:"code" db:"code" validate:"required,len=2,uppercase"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	SourceURL   *string   `json:"source_url,omitempty" db:"source_url" validate:"omitempty,url"`
	Attribution *string   `json:"attribution,omitempty" db:"attribution" validate:"omitempty,max=255"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateCountryRequest represents the request body for creating a country
type CreateCountryRequest struct {
	Code        string  `json:"code" validate:"required,len=2,uppercase"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	SourceURL   *string `json:"source_url,omitempty" validate:"omitempty,url"`
	Attribution *string `json:"attribution,omitempty" validate:"omitempty,max=255"`
}

// UpdateCountryRequest represents the request body for updating a country
type UpdateCountryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	SourceURL   *string `json:"source_url,omitempty" validate:"omitempty,url"`
	Attribution *string `json:"attribution,omitempty" validate:"omitempty,max=255"`
}

// CountryStats represents statistics for a country
type CountryStats struct {
	DatasetCount int `json:"dataset_count"`
	TotalNames   int `json:"total_names"`
	YearRange    struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"year_range"`
}

// CountryWithStats represents a country with its statistics
type CountryWithStats struct {
	Country
	Stats *CountryStats `json:"stats,omitempty"`
}
