package db

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type YearRange struct {
	MinYear int `json:"min_year"`
	MaxYear int `json:"max_year"`
}

func (db *DB) GetYearRange(ctx context.Context) (*YearRange, error) {
	query := `
		SELECT 
			MIN(year) as min_year,
			MAX(year) as max_year
		FROM names
	`

	var yr YearRange
	err := db.Pool.QueryRow(ctx, query).Scan(&yr.MinYear, &yr.MaxYear)
	if err != nil {
		return nil, err
	}

	return &yr, nil
}

type Country struct {
	Code                             string `json:"code"`
	Name                             string `json:"name"`
	DataSourceName                   string `json:"data_source_name"`
	DataSourceURL                    string `json:"data_source_url"`
	DataSourceDescription            string `json:"data_source_description"`
	DataSourceRequiresManualDownload bool   `json:"data_source_requires_manual_download"`
}

type CountriesResponse struct {
	Countries []Country `json:"countries"`
}

func (db *DB) GetCountries(ctx context.Context) (*CountriesResponse, error) {
	query := `
		SELECT 
			code,
			name,
			data_source_name,
			data_source_url,
			COALESCE(data_source_description, '') as data_source_description,
			data_source_requires_manual_download
		FROM countries
		ORDER BY name
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var c Country
		err := rows.Scan(
			&c.Code,
			&c.Name,
			&c.DataSourceName,
			&c.DataSourceURL,
			&c.DataSourceDescription,
			&c.DataSourceRequiresManualDownload,
		)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	return &CountriesResponse{Countries: countries}, nil
}

type NameRecord struct {
	Name            string   `json:"name"`
	TotalCount      int      `json:"total_count"`
	FemaleCount     int      `json:"female_count"`
	MaleCount       int      `json:"male_count"`
	GenderBalance   float64  `json:"gender_balance"`
	Rank            int      `json:"rank"`
	CumulativeShare float64  `json:"cumulative_share"`
	NameStart       int      `json:"name_start"`
	NameEnd         int      `json:"name_end"`
	Countries       []string `json:"countries"`
}

type NamesListMeta struct {
	Page              int                `json:"page"`
	PageSize          int                `json:"page_size"`
	TotalCount        int                `json:"total_count"`
	TotalPages        int                `json:"total_pages"`
	DbStart           int                `json:"db_start"`
	DbEnd             int                `json:"db_end"`
	PopularitySummary *PopularitySummary `json:"popularity_summary,omitempty"`
}

type PopularitySummary struct {
	PopulationTotal        int64   `json:"population_total"`
	ActiveDriver           string  `json:"active_driver"`
	ActiveValue            float64 `json:"active_value"`
	DerivedMinCount        int     `json:"derived_min_count"`
	DerivedTopN            int     `json:"derived_top_n"`
	DerivedCoveragePercent float64 `json:"derived_coverage_percent"`
}

type NamesListResponse struct {
	Meta  NamesListMeta `json:"meta"`
	Names []NameRecord  `json:"names"`
}

type NamesListParams struct {
	// Year filters
	YearFrom int
	YearTo   int

	// Country filter
	Countries []string

	// Gender balance filter (0-100)
	GenderBalanceMin int
	GenderBalanceMax int

	// Popularity filters (only one should be active)
	MinCount        int
	TopN            int
	CoveragePercent float64

	// Name pattern filter
	NameGlob string

	// Sorting
	SortKey   string // popularity, total_count, name, gender_balance, countries
	SortOrder string // asc, desc

	// Pagination
	Page     int
	PageSize int
}

func ParseNamesListParams(query url.Values, dbStart, dbEnd int) (*NamesListParams, error) {
	params := &NamesListParams{
		// Defaults
		YearFrom:         dbStart,
		YearTo:           dbEnd,
		Countries:        []string{}, // empty = all countries
		GenderBalanceMin: 0,
		GenderBalanceMax: 100,
		MinCount:         0,
		TopN:             0,
		CoveragePercent:  0,
		NameGlob:         "",
		SortKey:          "popularity",
		SortOrder:        "asc",
		Page:             1,
		PageSize:         50,
	}

	// Parse year_from
	if v := query.Get("year_from"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("year_from must be an integer")
		}
		params.YearFrom = val
	}

	// Parse year_to
	if v := query.Get("year_to"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("year_to must be an integer")
		}
		params.YearTo = val
	}

	// Parse countries (comma-separated)
	if v := query.Get("countries"); v != "" {
		params.Countries = strings.Split(v, ",")
	}

	// Parse gender_balance_min
	if v := query.Get("gender_balance_min"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("gender_balance_min must be an integer")
		}
		params.GenderBalanceMin = val
	}

	// Parse gender_balance_max
	if v := query.Get("gender_balance_max"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("gender_balance_max must be an integer")
		}
		params.GenderBalanceMax = val
	}

	// Parse min_count
	if v := query.Get("min_count"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("min_count must be an integer")
		}
		params.MinCount = val
	}

	// Parse top_n
	if v := query.Get("top_n"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("top_n must be an integer")
		}
		params.TopN = val
	}

	// Parse coverage_percent
	if v := query.Get("coverage_percent"); v != "" {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("coverage_percent must be a number")
		}
		params.CoveragePercent = val
	}

	// Parse name_glob
	if v := query.Get("name_glob"); v != "" {
		params.NameGlob = v
	}

	// Parse sort_key
	if v := query.Get("sort_key"); v != "" {
		params.SortKey = v
	}

	// Parse sort_order
	if v := query.Get("sort_order"); v != "" {
		params.SortOrder = v
	}

	// Parse page
	if v := query.Get("page"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("page must be an integer")
		}
		params.Page = val
	}

	// Parse page_size
	if v := query.Get("page_size"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("page_size must be an integer")
		}
		params.PageSize = val
	}

	// Validate
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}

func (p *NamesListParams) Validate() error {
	// Year range validation
	if p.YearFrom > p.YearTo {
		return fmt.Errorf("year_from must be <= year_to")
	}

	// Gender balance validation
	if p.GenderBalanceMin < 0 || p.GenderBalanceMin > 100 {
		return fmt.Errorf("gender_balance_min must be between 0 and 100")
	}
	if p.GenderBalanceMax < 0 || p.GenderBalanceMax > 100 {
		return fmt.Errorf("gender_balance_max must be between 0 and 100")
	}
	if p.GenderBalanceMin > p.GenderBalanceMax {
		return fmt.Errorf("gender_balance_min must be <= gender_balance_max")
	}

	// Page validation
	if p.Page < 1 || p.Page > 100 {
		return fmt.Errorf("page must be between 1 and 100")
	}

	// Page size validation
	if p.PageSize < 10 || p.PageSize > 100 {
		return fmt.Errorf("page_size must be between 10 and 100")
	}

	// Sort key validation
	validSortKeys := map[string]bool{
		"popularity":     true,
		"total_count":    true,
		"name":           true,
		"gender_balance": true,
		"countries":      true,
	}
	if !validSortKeys[p.SortKey] {
		return fmt.Errorf("sort_key must be one of: popularity, total_count, name, gender_balance, countries")
	}

	// Sort order validation
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		return fmt.Errorf("sort_order must be either 'asc' or 'desc'")
	}

	return nil
}

// GetActivePopularityFilter returns which popularity filter is active
func (p *NamesListParams) GetActivePopularityFilter() string {
	if p.CoveragePercent > 0 {
		return "coverage_percent"
	}
	if p.TopN > 0 {
		return "top_n"
	}
	if p.MinCount > 0 {
		return "min_count"
	}
	return "none"
}

func (db *DB) GetNamesList(ctx context.Context, params *NamesListParams) (*NamesListResponse, error) {
	// This is a complex 6-stage pipeline query
	// For simplicity, we'll build it with CTEs (Common Table Expressions)

	query := `
	WITH 
	-- Stage 1: Basic Filters
	filtered_names AS (
		SELECT 
			n.name,
			n.year,
			n.gender,
			n.count,
			c.code as country_code
		FROM names n
		JOIN countries c ON n.country_id = c.id
		WHERE n.year >= $1 
		  AND n.year <= $2
		  AND ($3::text[] IS NULL OR c.code = ANY($3::text[]))
		  AND ($4 = '' OR n.name ILIKE $4)
	),
	-- Stage 2: Aggregation
	aggregated AS (
		SELECT 
			name,
			SUM(count) as total_count,
			SUM(CASE WHEN gender = 'F' THEN count ELSE 0 END) as female_count,
			SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END) as male_count,
			CASE 
				WHEN SUM(CASE WHEN gender IN ('M','F') THEN count ELSE 0 END) = 0 THEN NULL
				ELSE 100.0 * SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END)::float / 
				     NULLIF(SUM(CASE WHEN gender IN ('M','F') THEN count ELSE 0 END), 0)
			END as gender_balance,
			MIN(year) as name_start,
			MAX(year) as name_end,
			ARRAY_AGG(DISTINCT country_code ORDER BY country_code) as countries
		FROM filtered_names
		GROUP BY name
	),
	-- Stage 3: Gender Balance Filter
	gender_filtered AS (
		SELECT *
		FROM aggregated
		WHERE (gender_balance IS NULL OR (gender_balance >= $5 AND gender_balance <= $6))
	),
	-- Stage 4: Popularity Computation
	ranked AS (
		SELECT
			*,
			ROW_NUMBER() OVER (ORDER BY total_count DESC, name ASC) as rank,
			SUM(total_count) OVER (ORDER BY total_count DESC, name ASC) as cumulative_count,
			SUM(total_count) OVER () as population_total
		FROM gender_filtered
	),
	with_cumulative_share AS (
		SELECT
			*,
			cumulative_count::float / NULLIF(population_total, 0) as cumulative_share
		FROM ranked
	),
	-- Stage 5: Popularity Filter
	popularity_filtered AS (
		SELECT *
		FROM with_cumulative_share
		WHERE
			CASE
				WHEN $7 > 0 THEN cumulative_share <= $7 / 100.0  -- coverage_percent
				WHEN $8 > 0 THEN rank <= $8                       -- top_n
				WHEN $9 > 0 THEN total_count >= $9                -- min_count
				ELSE true                                          -- no filter
			END
	),
	-- Count total for pagination
	total_count AS (
		SELECT COUNT(*) as cnt FROM popularity_filtered
	)
	-- Stage 6: Final Sorting & Pagination
	SELECT 
		pf.*,
		tc.cnt as total_count_val
	FROM popularity_filtered pf
	CROSS JOIN total_count tc
	ORDER BY 
		CASE WHEN $10 = 'popularity' AND $11 = 'asc' THEN pf.rank END ASC,
		CASE WHEN $10 = 'popularity' AND $11 = 'desc' THEN pf.rank END DESC,
		CASE WHEN $10 = 'total_count' AND $11 = 'asc' THEN pf.total_count END ASC,
		CASE WHEN $10 = 'total_count' AND $11 = 'desc' THEN pf.total_count END DESC,
		CASE WHEN $10 = 'name' AND $11 = 'asc' THEN pf.name END ASC,
		CASE WHEN $10 = 'name' AND $11 = 'desc' THEN pf.name END DESC,
		CASE WHEN $10 = 'gender_balance' AND $11 = 'asc' THEN pf.gender_balance END ASC,
		CASE WHEN $10 = 'gender_balance' AND $11 = 'desc' THEN pf.gender_balance END DESC,
		-- Tie breakers
		pf.total_count DESC,
		pf.name ASC
	LIMIT $12 OFFSET $13
	`

	// Convert name_glob to SQL ILIKE pattern
	globPattern := ""
	if params.NameGlob != "" {
		// Convert * to % and ? to _ for SQL LIKE
		globPattern = strings.ReplaceAll(params.NameGlob, "*", "%")
		globPattern = strings.ReplaceAll(globPattern, "?", "_")
	}

	// Handle country filter (nil for all countries)
	var countries interface{}
	if len(params.Countries) == 0 {
		countries = nil
	} else {
		countries = params.Countries
	}

	offset := (params.Page - 1) * params.PageSize

	rows, err := db.Pool.Query(ctx, query,
		params.YearFrom,         // $1
		params.YearTo,           // $2
		countries,               // $3
		globPattern,             // $4
		params.GenderBalanceMin, // $5
		params.GenderBalanceMax, // $6
		params.CoveragePercent,  // $7
		params.TopN,             // $8
		params.MinCount,         // $9
		params.SortKey,          // $10
		params.SortOrder,        // $11
		params.PageSize,         // $12
		offset,                  // $13
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var names []NameRecord
	var totalCount int
	var populationTotal int64

	for rows.Next() {
		var nr NameRecord
		var genderBalance *float64
		var totalCountVal int
		var cumulativeCount int64

		err := rows.Scan(
			&nr.Name,
			&nr.TotalCount,
			&nr.FemaleCount,
			&nr.MaleCount,
			&genderBalance,
			&nr.NameStart,
			&nr.NameEnd,
			&nr.Countries,
			&nr.Rank,
			&cumulativeCount,
			&populationTotal,
			&nr.CumulativeShare,
			&totalCountVal,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if genderBalance != nil {
			nr.GenderBalance = *genderBalance
		}

		totalCount = totalCountVal
		names = append(names, nr)
	}

	// Get database year range for meta
	yearRange, _ := db.GetYearRange(ctx)

	// Calculate total pages
	totalPages := (totalCount + params.PageSize - 1) / params.PageSize

	meta := NamesListMeta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		DbStart:    yearRange.MinYear,
		DbEnd:      yearRange.MaxYear,
	}

	// Add popularity summary if filter is active
	activeDriver := params.GetActivePopularityFilter()
	if activeDriver != "none" && len(names) > 0 {
		meta.PopularitySummary = &PopularitySummary{
			PopulationTotal: populationTotal,
			ActiveDriver:    activeDriver,
		}

		// Set active value and derived values based on driver
		switch activeDriver {
		case "coverage_percent":
			meta.PopularitySummary.ActiveValue = params.CoveragePercent
			meta.PopularitySummary.DerivedCoveragePercent = params.CoveragePercent
			if len(names) > 0 {
				meta.PopularitySummary.DerivedTopN = names[len(names)-1].Rank
				meta.PopularitySummary.DerivedMinCount = names[len(names)-1].TotalCount
			}
		case "top_n":
			meta.PopularitySummary.ActiveValue = float64(params.TopN)
			meta.PopularitySummary.DerivedTopN = params.TopN
			if len(names) > 0 {
				meta.PopularitySummary.DerivedMinCount = names[len(names)-1].TotalCount
				meta.PopularitySummary.DerivedCoveragePercent = names[len(names)-1].CumulativeShare * 100
			}
		case "min_count":
			meta.PopularitySummary.ActiveValue = float64(params.MinCount)
			meta.PopularitySummary.DerivedMinCount = params.MinCount
			if len(names) > 0 {
				meta.PopularitySummary.DerivedTopN = names[len(names)-1].Rank
				meta.PopularitySummary.DerivedCoveragePercent = names[len(names)-1].CumulativeShare * 100
			}
		}
	}

	return &NamesListResponse{
		Meta:  meta,
		Names: names,
	}, nil
}
