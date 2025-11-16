package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/db"
)

// NamesList returns the list of names
func NamesList(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.FixtureMode {
			// Load and return fixture JSON
			data, err := LoadFixture("../spec-examples/names-list.json")
			if err != nil {
				WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			WriteJSON(w, http.StatusOK, data)
			return
		}

		// Get year range for defaults
		ctx := r.Context()
		yearRange, err := cfg.DB.GetYearRange(ctx)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// Parse and validate parameters
		params, err := db.ParseNamesListParams(r.URL.Query(), yearRange.MinYear, yearRange.MaxYear)
		if err != nil {
			WriteError(w, http.StatusBadRequest, fmt.Sprintf("Invalid parameters: %v", err))
			return
		}

		// Query database
		response, err := cfg.DB.GetNamesList(ctx, params)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// NameTrend returns the trend details for a specific name
func NameTrend(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.FixtureMode {
			// Load and return fixture JSON
			data, err := LoadFixture("../spec-examples/name-detail.json")
			if err != nil {
				WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			WriteJSON(w, http.StatusOK, data)
			return
		}

		// Real database logic (stub for Phase 2)
		WriteError(w, http.StatusNotImplemented, "Database mode not implemented")
	}
}
