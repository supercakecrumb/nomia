package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/supercakecrumb/affirm-name/internal/model"
)

var (
	// ErrCountryNotFound is returned when a country is not found
	ErrCountryNotFound = errors.New("country not found")
	// ErrCountryCodeExists is returned when a country code already exists
	ErrCountryCodeExists = errors.New("country code already exists")
)

// CountryRepository handles database operations for countries
type CountryRepository struct {
	db *pgxpool.Pool
}

// NewCountryRepository creates a new country repository
func NewCountryRepository(db *pgxpool.Pool) *CountryRepository {
	return &CountryRepository{db: db}
}

// Create creates a new country
func (r *CountryRepository) Create(ctx context.Context, country *model.Country) error {
	query := `
		INSERT INTO countries (id, code, name, source_url, attribution, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		country.ID,
		country.Code,
		country.Name,
		country.SourceURL,
		country.Attribution,
	).Scan(&country.ID, &country.CreatedAt, &country.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation on code
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"countries_code_key\" (SQLSTATE 23505)" {
			return ErrCountryCodeExists
		}
		return fmt.Errorf("failed to create country: %w", err)
	}

	return nil
}

// GetByID retrieves a country by its ID
func (r *CountryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Country, error) {
	query := `
		SELECT id, code, name, source_url, attribution, created_at, updated_at
		FROM countries
		WHERE id = $1
	`

	var country model.Country
	err := r.db.QueryRow(ctx, query, id).Scan(
		&country.ID,
		&country.Code,
		&country.Name,
		&country.SourceURL,
		&country.Attribution,
		&country.CreatedAt,
		&country.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to get country by id: %w", err)
	}

	return &country, nil
}

// GetByCode retrieves a country by its code
func (r *CountryRepository) GetByCode(ctx context.Context, code string) (*model.Country, error) {
	query := `
		SELECT id, code, name, source_url, attribution, created_at, updated_at
		FROM countries
		WHERE code = $1
	`

	var country model.Country
	err := r.db.QueryRow(ctx, query, code).Scan(
		&country.ID,
		&country.Code,
		&country.Name,
		&country.SourceURL,
		&country.Attribution,
		&country.CreatedAt,
		&country.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCountryNotFound
		}
		return nil, fmt.Errorf("failed to get country by code: %w", err)
	}

	return &country, nil
}

// List retrieves a paginated list of countries
func (r *CountryRepository) List(ctx context.Context, limit, offset int) ([]*model.Country, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM countries`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count countries: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, code, name, source_url, attribution, created_at, updated_at
		FROM countries
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list countries: %w", err)
	}
	defer rows.Close()

	var countries []*model.Country
	for rows.Next() {
		var country model.Country
		err := rows.Scan(
			&country.ID,
			&country.Code,
			&country.Name,
			&country.SourceURL,
			&country.Attribution,
			&country.CreatedAt,
			&country.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, &country)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating countries: %w", err)
	}

	return countries, total, nil
}

// Update updates an existing country
func (r *CountryRepository) Update(ctx context.Context, country *model.Country) error {
	query := `
		UPDATE countries
		SET name = $1, source_url = $2, attribution = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		country.Name,
		country.SourceURL,
		country.Attribution,
		country.ID,
	).Scan(&country.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCountryNotFound
		}
		return fmt.Errorf("failed to update country: %w", err)
	}

	return nil
}

// Delete deletes a country by its ID
func (r *CountryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM countries WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		// Check for foreign key constraint violation
		if err.Error() == "ERROR: update or delete on table \"countries\" violates foreign key constraint \"datasets_country_id_fkey\" on table \"datasets\" (SQLSTATE 23503)" {
			return errors.New("cannot delete country with associated datasets")
		}
		return fmt.Errorf("failed to delete country: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCountryNotFound
	}

	return nil
}

// GetStats retrieves statistics for a country
func (r *CountryRepository) GetStats(ctx context.Context, countryID uuid.UUID) (*model.CountryStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT d.id) as dataset_count,
			COUNT(n.id) as total_names,
			COALESCE(MIN(n.year), 0) as min_year,
			COALESCE(MAX(n.year), 0) as max_year
		FROM countries c
		LEFT JOIN datasets d ON d.country_id = c.id AND d.deleted_at IS NULL AND d.status = 'completed'
		LEFT JOIN names n ON n.country_id = c.id AND n.deleted_at IS NULL
		WHERE c.id = $1
		GROUP BY c.id
	`

	var stats model.CountryStats
	err := r.db.QueryRow(ctx, query, countryID).Scan(
		&stats.DatasetCount,
		&stats.TotalNames,
		&stats.YearRange.Min,
		&stats.YearRange.Max,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Country exists but has no data
			return &model.CountryStats{}, nil
		}
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}

	return &stats, nil
}
