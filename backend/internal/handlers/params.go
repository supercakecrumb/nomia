package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type NameTrendParams struct {
	Name      string   // Required
	YearFrom  int      // Optional, defaults to db_start
	YearTo    int      // Optional, defaults to db_end
	Countries []string // Optional, defaults to all countries
}

func ParseNameTrendParams(query url.Values, dbStart, dbEnd int) (*NameTrendParams, error) {
	params := &NameTrendParams{
		YearFrom:  dbStart,
		YearTo:    dbEnd,
		Countries: []string{}, // empty = all countries
	}

	// Parse name (required)
	name := query.Get("name")
	if name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}
	params.Name = name

	// Parse year_min (frontend parameter name)
	if v := query.Get("year_min"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("year_min must be an integer")
		}
		params.YearFrom = val
	}

	// Parse year_max (frontend parameter name)
	if v := query.Get("year_max"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("year_max must be an integer")
		}
		params.YearTo = val
	}

	// Parse countries (comma-separated)
	if v := query.Get("countries"); v != "" {
		params.Countries = strings.Split(v, ",")
	}

	// Validate
	if params.YearFrom > params.YearTo {
		return nil, fmt.Errorf("year_from must be <= year_to")
	}

	return params, nil
}
