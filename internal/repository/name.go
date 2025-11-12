package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/supercakecrumb/affirm-name/internal/model"
)

// NameRepository handles database operations for name records
type NameRepository struct {
	db *sqlx.DB
}

// NewNameRepository creates a new name repository
func NewNameRepository(db *sqlx.DB) *NameRepository {
	return &NameRepository{db: db}
}

// BatchInsert inserts name records using staging table pattern for atomic operations
// Returns the number of rows inserted
func (r *NameRepository) BatchInsert(
	ctx context.Context,
	datasetID uuid.UUID,
	countryID uuid.UUID,
	year int,
	records []*model.NameRecord,
) (int, error) {
	if len(records) == 0 {
		return 0, nil
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create temporary staging table
	stagingTable := "names_staging"
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		CREATE TEMPORARY TABLE %s (
			dataset_id UUID NOT NULL,
			country_id UUID NOT NULL,
			year INTEGER NOT NULL,
			name VARCHAR(100) NOT NULL,
			gender CHAR(1) NOT NULL,
			count INTEGER NOT NULL
		) ON COMMIT DROP
	`, stagingTable))
	if err != nil {
		return 0, fmt.Errorf("failed to create staging table: %w", err)
	}

	// Batch insert into staging table
	const batchSize = 1000
	totalInserted := 0

	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		// Build bulk insert query
		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*6)

		for _, record := range batch {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				datasetID,
				countryID,
				year,
				record.Name,
				record.Gender,
				record.Count,
			)
		}

		query := fmt.Sprintf(`
			INSERT INTO %s (dataset_id, country_id, year, name, gender, count)
			VALUES %s
		`, stagingTable, strings.Join(valueStrings, ","))

		// Rebind for PostgreSQL
		query = tx.Rebind(query)

		result, err := tx.ExecContext(ctx, query, valueArgs...)
		if err != nil {
			return 0, fmt.Errorf("failed to insert batch into staging: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to get rows affected: %w", err)
		}
		totalInserted += int(rows)
	}

	// Validate staging data
	var invalidCount int
	err = tx.GetContext(ctx, &invalidCount, fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE name = '' OR gender NOT IN ('M', 'F') OR count < 0
	`, stagingTable))
	if err != nil {
		return 0, fmt.Errorf("failed to validate staging data: %w", err)
	}
	if invalidCount > 0 {
		return 0, fmt.Errorf("staging data contains %d invalid records", invalidCount)
	}

	// Move data from staging to main table
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO names (dataset_id, country_id, year, name, gender, count)
		SELECT dataset_id, country_id, year, name, gender, count
		FROM %s
	`, stagingTable))
	if err != nil {
		return 0, fmt.Errorf("failed to move data from staging to main table: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return totalInserted, nil
}

// ListWithRank retrieves names with optional filters, pagination, and rank calculation
func (r *NameRepository) ListWithRank(
	ctx context.Context,
	filters *model.NameFilters,
	sortBy string, // "count" or "name"
	sortOrder string, // "asc" or "desc"
	limit, offset int,
) ([]*model.Name, int, error) {
	// Build WHERE clause
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if filters != nil {
		if filters.DatasetID != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.dataset_id = $%d", argPos))
			args = append(args, *filters.DatasetID)
			argPos++
		}
		if filters.CountryID != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.country_id = $%d", argPos))
			args = append(args, *filters.CountryID)
			argPos++
		}
		if filters.Year != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.year = $%d", argPos))
			args = append(args, *filters.Year)
			argPos++
		}
		if filters.Name != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.name ILIKE $%d", argPos))
			args = append(args, *filters.Name+"%")
			argPos++
		}
		if filters.Gender != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.gender = $%d", argPos))
			args = append(args, *filters.Gender)
			argPos++
		}
		if filters.MinCount != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("n.count >= $%d", argPos))
			args = append(args, *filters.MinCount)
			argPos++
		}
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM names n %s", whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count names: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "n.count DESC, n.name ASC"
	if sortBy == "name" {
		if sortOrder == "desc" {
			orderBy = "n.name DESC"
		} else {
			orderBy = "n.name ASC"
		}
	} else if sortBy == "count" {
		if sortOrder == "asc" {
			orderBy = "n.count ASC, n.name ASC"
		} else {
			orderBy = "n.count DESC, n.name ASC"
		}
	}

	// Get paginated results with rank
	query := fmt.Sprintf(`
		SELECT
			n.id,
			n.dataset_id,
			n.country_id,
			n.year,
			n.name,
			n.gender,
			n.count,
			n.created_at
		FROM names n
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argPos, argPos+1)

	args = append(args, limit, offset)

	var names []*model.Name
	err = r.db.SelectContext(ctx, &names, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list names: %w", err)
	}

	return names, total, nil
}

// List retrieves names with optional filters, pagination (legacy method)
func (r *NameRepository) List(
	ctx context.Context,
	filters *model.NameFilters,
	limit, offset int,
) ([]*model.Name, int, error) {
	return r.ListWithRank(ctx, filters, "count", "desc", limit, offset)
}

// Search performs prefix search across years and returns aggregated results
func (r *NameRepository) Search(
	ctx context.Context,
	query string,
	countryID *uuid.UUID,
	gender *string,
	limit, offset int,
) ([]*model.NameSearchResult, int, error) {
	whereClauses := []string{"n.name ILIKE $1"}
	args := []interface{}{query + "%"}
	argPos := 2

	if countryID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("n.country_id = $%d", argPos))
		args = append(args, *countryID)
		argPos++
	}
	if gender != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("n.gender = $%d", argPos))
		args = append(args, *gender)
		argPos++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count of unique names
	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT n.name)
		FROM names n
		WHERE %s
	`, whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get aggregated search results
	searchQuery := fmt.Sprintf(`
		SELECT
			n.name,
			SUM(n.count) as total_count,
			ARRAY_AGG(DISTINCT c.code ORDER BY c.code) as countries,
			MIN(n.year) as min_year,
			MAX(n.year) as max_year,
			CASE
				WHEN COUNT(DISTINCT n.gender) = 1 THEN MAX(n.gender)
				ELSE NULL
			END as primary_gender,
			SUM(CASE WHEN n.gender = 'M' THEN n.count ELSE 0 END)::FLOAT / NULLIF(SUM(n.count), 0) * 100 as male_percentage,
			SUM(CASE WHEN n.gender = 'F' THEN n.count ELSE 0 END)::FLOAT / NULLIF(SUM(n.count), 0) * 100 as female_percentage
		FROM names n
		JOIN countries c ON n.country_id = c.id
		WHERE %s
		GROUP BY n.name
		ORDER BY SUM(n.count) DESC, n.name ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)

	var results []*model.NameSearchResult
	err = r.db.SelectContext(ctx, &results, searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search names: %w", err)
	}

	return results, total, nil
}

// GetByName retrieves all records for a specific name with optional filters
func (r *NameRepository) GetByName(
	ctx context.Context,
	name string,
	filters *model.NameFilters,
) ([]*model.Name, error) {
	whereClauses := []string{"name = $1"}
	args := []interface{}{name}
	argPos := 2

	if filters != nil {
		if filters.DatasetID != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("dataset_id = $%d", argPos))
			args = append(args, *filters.DatasetID)
			argPos++
		}
		if filters.CountryID != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("country_id = $%d", argPos))
			args = append(args, *filters.CountryID)
			argPos++
		}
		if filters.Year != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("year = $%d", argPos))
			args = append(args, *filters.Year)
			argPos++
		}
		if filters.Gender != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("gender = $%d", argPos))
			args = append(args, *filters.Gender)
			argPos++
		}
	}

	query := fmt.Sprintf(`
		SELECT id, dataset_id, country_id, year, name, gender, count, created_at
		FROM names
		WHERE %s
		ORDER BY year DESC, count DESC
	`, strings.Join(whereClauses, " AND "))

	var names []*model.Name
	err := r.db.SelectContext(ctx, &names, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get name records: %w", err)
	}

	return names, nil
}

// CalculateRanks calculates ranks for names within year/country/gender groups
func (r *NameRepository) CalculateRanks(
	ctx context.Context,
	countryID uuid.UUID,
	year int,
	gender *string,
) (map[string]int, error) {
	whereClauses := []string{
		"country_id = $1",
		"year = $2",
	}
	args := []interface{}{countryID, year}
	argPos := 3

	if gender != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, *gender)
		argPos++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
		SELECT
			name,
			ROW_NUMBER() OVER (ORDER BY count DESC, name ASC) as rank
		FROM names
		WHERE %s
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate ranks: %w", err)
	}
	defer rows.Close()

	ranks := make(map[string]int)
	for rows.Next() {
		var name string
		var rank int
		if err := rows.Scan(&name, &rank); err != nil {
			return nil, fmt.Errorf("failed to scan rank: %w", err)
		}
		ranks[name] = rank
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ranks: %w", err)
	}

	return ranks, nil
}

// Delete removes all name records for a dataset (used for reprocessing)
func (r *NameRepository) Delete(ctx context.Context, datasetID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM names WHERE dataset_id = $1
	`, datasetID)
	if err != nil {
		return fmt.Errorf("failed to delete names: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetStats returns statistics for a dataset
func (r *NameRepository) GetStats(ctx context.Context, datasetID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total records
	var total int
	err := r.db.GetContext(ctx, &total, `
		SELECT COUNT(*) FROM names WHERE dataset_id = $1
	`, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	stats["total_records"] = total

	// Unique names
	var uniqueNames int
	err = r.db.GetContext(ctx, &uniqueNames, `
		SELECT COUNT(DISTINCT name) FROM names WHERE dataset_id = $1
	`, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique names: %w", err)
	}
	stats["unique_names"] = uniqueNames

	// Year range
	var minYear, maxYear sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT MIN(year), MAX(year) FROM names WHERE dataset_id = $1
	`, datasetID).Scan(&minYear, &maxYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get year range: %w", err)
	}
	if minYear.Valid {
		stats["min_year"] = minYear.Int64
	}
	if maxYear.Valid {
		stats["max_year"] = maxYear.Int64
	}

	return stats, nil
}
