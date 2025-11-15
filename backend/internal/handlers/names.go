package handlers

import (
	"net/http"

	"github.com/supercakecrumb/affirm-name/internal/config"
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

		// Real database logic (stub for Phase 2)
		WriteError(w, http.StatusNotImplemented, "Database mode not implemented")
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
