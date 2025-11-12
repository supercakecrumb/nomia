package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supercakecrumb/affirm-name/internal/api/handlers"
	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/database"
	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/service"
)

func setupNameTestRouter(t *testing.T) (*gin.Engine, *database.DB, func()) {
	// Load test configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Connect to test database
	db, err := database.New(cfg.Database)
	require.NoError(t, err)

	// Initialize repositories
	countryRepo := repository.NewCountryRepository(db.Pool)
	nameRepo := repository.NewNameRepository(db.Sqlx)

	// Initialize services
	nameService := service.NewNameService(nameRepo, countryRepo)

	// Initialize handlers
	nameHandler := handlers.NewNameHandler(nameService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/v1")
	{
		names := v1.Group("/names")
		{
			names.GET("", nameHandler.List)
			names.GET("/search", nameHandler.Search)
			names.GET("/:name", nameHandler.GetByName)
		}
	}

	cleanup := func() {
		db.Close()
	}

	return router, db, cleanup
}

func TestNameList(t *testing.T) {
	router, db, cleanup := setupNameTestRouter(t)
	defer cleanup()

	ctx := context.Background()

	// Create test country
	countryRepo := repository.NewCountryRepository(db.Pool)
	country := &model.Country{
		ID:   uuid.New(),
		Code: "US",
		Name: "United States",
	}
	err := countryRepo.Create(ctx, country)
	require.NoError(t, err)

	// Create test dataset
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	datasetReq := &model.CreateDatasetRequest{
		CountryID:  country.ID,
		Filename:   "test.csv",
		FilePath:   "/tmp/test.csv",
		FileSize:   1000,
		UploadedBy: "test@example.com",
	}
	dataset, err := datasetRepo.Create(ctx, datasetReq)
	require.NoError(t, err)

	// Insert test names
	nameRepo := repository.NewNameRepository(db.Sqlx)
	records := []*model.NameRecord{
		{Year: 2020, Name: "Emma", Gender: "F", Count: 15581},
		{Year: 2020, Name: "Olivia", Gender: "F", Count: 17535},
		{Year: 2020, Name: "Liam", Gender: "M", Count: 19659},
		{Year: 2020, Name: "Noah", Gender: "M", Count: 18252},
	}
	_, err = nameRepo.BatchInsert(ctx, dataset.ID, country.ID, 2020, records)
	require.NoError(t, err)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "List names with required parameters",
			queryParams:    "country=US&year=2020",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 4)

				meta := body["meta"].(map[string]interface{})
				assert.Equal(t, float64(4), meta["total"])
			},
		},
		{
			name:           "List names with gender filter",
			queryParams:    "country=US&year=2020&gender=F",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.Equal(t, 2, len(data))

				// Check all results are female
				for _, item := range data {
					name := item.(map[string]interface{})
					assert.Equal(t, "F", name["gender"])
				}
			},
		},
		{
			name:           "List names with name prefix",
			queryParams:    "country=US&year=2020&name=Em",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.Equal(t, 1, len(data))

				name := data[0].(map[string]interface{})
				assert.Equal(t, "Emma", name["name"])
			},
		},
		{
			name:           "List names with sorting by name",
			queryParams:    "country=US&year=2020&sort=name:asc",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 4)

				// Check first name is alphabetically first
				firstName := data[0].(map[string]interface{})
				assert.Equal(t, "Emma", firstName["name"])
			},
		},
		{
			name:           "List names with pagination",
			queryParams:    "country=US&year=2020&limit=2&offset=0",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.Equal(t, 2, len(data))

				meta := body["meta"].(map[string]interface{})
				assert.Equal(t, float64(4), meta["total"])
				assert.Equal(t, true, meta["has_more"])
			},
		},
		{
			name:           "Missing country parameter",
			queryParams:    "year=2020",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
		{
			name:           "Missing year parameter",
			queryParams:    "country=US",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
		{
			name:           "Invalid year",
			queryParams:    "country=US&year=1800",
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
		{
			name:           "Invalid gender",
			queryParams:    "country=US&year=2020&gender=X",
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
		{
			name:           "Country not found",
			queryParams:    "country=XX&year=2020",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "not_found", errorObj["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/names?"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestNameSearch(t *testing.T) {
	router, db, cleanup := setupNameTestRouter(t)
	defer cleanup()

	ctx := context.Background()

	// Create test country
	countryRepo := repository.NewCountryRepository(db.Pool)
	country := &model.Country{
		ID:   uuid.New(),
		Code: "US",
		Name: "United States",
	}
	err := countryRepo.Create(ctx, country)
	require.NoError(t, err)

	// Create test dataset
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	datasetReq := &model.CreateDatasetRequest{
		CountryID:  country.ID,
		Filename:   "test.csv",
		FilePath:   "/tmp/test.csv",
		FileSize:   1000,
		UploadedBy: "test@example.com",
	}
	dataset, err := datasetRepo.Create(ctx, datasetReq)
	require.NoError(t, err)

	// Insert test names across multiple years
	nameRepo := repository.NewNameRepository(db.Sqlx)
	records2019 := []*model.NameRecord{
		{Year: 2019, Name: "Emma", Gender: "F", Count: 15000},
	}
	_, err = nameRepo.BatchInsert(ctx, dataset.ID, country.ID, 2019, records2019)
	require.NoError(t, err)

	records2020 := []*model.NameRecord{
		{Year: 2020, Name: "Emma", Gender: "F", Count: 15581},
		{Year: 2020, Name: "Emily", Gender: "F", Count: 12000},
	}
	_, err = nameRepo.BatchInsert(ctx, dataset.ID, country.ID, 2020, records2020)
	require.NoError(t, err)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "Search with valid query",
			queryParams:    "q=Em",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 2)
			},
		},
		{
			name:           "Search with country filter",
			queryParams:    "q=Em&country=US",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 2)
			},
		},
		{
			name:           "Search with gender filter",
			queryParams:    "q=Em&gender=F",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 2)
			},
		},
		{
			name:           "Missing query parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
		{
			name:           "Query too short",
			queryParams:    "q=E",
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "validation_error", errorObj["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/names/search?"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestNameGetByName(t *testing.T) {
	router, db, cleanup := setupNameTestRouter(t)
	defer cleanup()

	ctx := context.Background()

	// Create test country
	countryRepo := repository.NewCountryRepository(db.Pool)
	country := &model.Country{
		ID:   uuid.New(),
		Code: "US",
		Name: "United States",
	}
	err := countryRepo.Create(ctx, country)
	require.NoError(t, err)

	// Create test dataset
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	datasetReq := &model.CreateDatasetRequest{
		CountryID:  country.ID,
		Filename:   "test.csv",
		FilePath:   "/tmp/test.csv",
		FileSize:   1000,
		UploadedBy: "test@example.com",
	}
	dataset, err := datasetRepo.Create(ctx, datasetReq)
	require.NoError(t, err)

	// Insert test names
	nameRepo := repository.NewNameRepository(db.Sqlx)
	records := []*model.NameRecord{
		{Year: 2020, Name: "Emma", Gender: "F", Count: 15581},
	}
	_, err = nameRepo.BatchInsert(ctx, dataset.ID, country.ID, 2020, records)
	require.NoError(t, err)

	tests := []struct {
		name           string
		namePath       string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "Get existing name",
			namePath:       "Emma",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Emma", data["name"])
				assert.Greater(t, data["total_count"], float64(0))
			},
		},
		{
			name:           "Get name with country filter",
			namePath:       "Emma",
			queryParams:    "country=US",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Emma", data["name"])
			},
		},
		{
			name:           "Name not found",
			namePath:       "NonExistentName",
			queryParams:    "",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				errorObj := body["error"].(map[string]interface{})
				assert.Equal(t, "not_found", errorObj["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/v1/names/%s", tt.namePath)
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}
