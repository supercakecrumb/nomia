package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/supercakecrumb/affirm-name/internal/config"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
}

func Health(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "not_connected"
		if cfg.DB != nil {
			// Test database connection
			ctx := r.Context()
			err := cfg.DB.Pool.Ping(ctx)
			if err == nil {
				dbStatus = "connected"
			} else {
				dbStatus = "error"
			}
		}

		response := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			Version:   "1.0.0",
			Database:  dbStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
