package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supercakecrumb/affirm-name/internal/config"
)

// MetaYears returns the list of available years
func MetaYears(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.FixtureMode {
			// Load and return fixture JSON
			data, err := LoadFixture("../spec-examples/meta-years.json")
			if err != nil {
				WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			WriteJSON(w, http.StatusOK, data)
			return
		}

		// Query database
		ctx := r.Context()
		yearRange, err := cfg.DB.GetYearRange(ctx)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(yearRange)
	}
}

// MetaCountries returns the list of available countries
func MetaCountries(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.FixtureMode {
			// Load and return fixture JSON
			data, err := LoadFixture("../spec-examples/countries.json")
			if err != nil {
				WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			WriteJSON(w, http.StatusOK, data)
			return
		}

		// Query database
		ctx := r.Context()
		countries, err := cfg.DB.GetCountries(ctx)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countries)
	}
}
