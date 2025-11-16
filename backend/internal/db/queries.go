package db

import (
	"context"
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
