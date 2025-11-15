package handlers

import (
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

		// Real database logic (stub for Phase 2)
		WriteError(w, http.StatusNotImplemented, "Database mode not implemented")
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

		// Real database logic (stub for Phase 2)
		WriteError(w, http.StatusNotImplemented, "Database mode not implemented")
	}
}
