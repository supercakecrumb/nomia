package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supercakecrumb/affirm-name/internal/service"
)

// NameHandler handles HTTP requests for name operations
type NameHandler struct {
	service *service.NameService
}

// NewNameHandler creates a new name handler
func NewNameHandler(service *service.NameService) *NameHandler {
	return &NameHandler{service: service}
}

// List handles GET /v1/names
func (h *NameHandler) List(c *gin.Context) {
	// Parse required parameters
	country := c.Query("country")
	yearStr := c.Query("year")

	if country == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "country parameter is required",
				Details: []FieldError{{Field: "country", Message: "required"}},
			},
		})
		return
	}

	if yearStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "year parameter is required",
				Details: []FieldError{{Field: "year", Message: "required"}},
			},
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1970 || year > 2030 {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "year must be between 1970 and 2030",
				Details: []FieldError{{Field: "year", Message: "must be between 1970 and 2030"}},
			},
		})
		return
	}

	// Parse optional parameters
	var gender *string
	if g := c.Query("gender"); g != "" {
		if g != "M" && g != "F" {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: "gender must be M or F",
					Details: []FieldError{{Field: "gender", Message: "must be M or F"}},
				},
			})
			return
		}
		gender = &g
	}

	var namePrefix *string
	if n := c.Query("name"); n != "" {
		if len(n) < 2 {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: "name prefix must be at least 2 characters",
					Details: []FieldError{{Field: "name", Message: "must be at least 2 characters"}},
				},
			})
			return
		}
		namePrefix = &n
	}

	var minCount *int
	if mc := c.Query("min_count"); mc != "" {
		count, err := strconv.Atoi(mc)
		if err != nil || count < 0 {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: "min_count must be a positive integer",
					Details: []FieldError{{Field: "min_count", Message: "must be a positive integer"}},
				},
			})
			return
		}
		minCount = &count
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit < 1 || limit > 1000 {
		limit = 100
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if offset < 0 {
		offset = 0
	}

	// Parse sort parameters
	sortBy := "count"
	sortOrder := "desc"
	if sort := c.Query("sort"); sort != "" {
		// Parse sort format: "field:direction"
		// e.g., "count:desc" or "name:asc"
		if sort == "name:asc" {
			sortBy = "name"
			sortOrder = "asc"
		} else if sort == "name:desc" {
			sortBy = "name"
			sortOrder = "desc"
		} else if sort == "count:asc" {
			sortBy = "count"
			sortOrder = "asc"
		} else if sort == "count:desc" {
			sortBy = "count"
			sortOrder = "desc"
		}
	}

	// Call service
	names, total, err := h.service.List(
		c.Request.Context(),
		country,
		year,
		gender,
		namePrefix,
		minCount,
		sortBy,
		sortOrder,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, service.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "not_found",
					Message: "Country not found",
				},
			})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to list names",
			},
		})
		return
	}

	// Build filters metadata
	filters := map[string]interface{}{
		"country": country,
		"year":    year,
	}
	if gender != nil {
		filters["gender"] = *gender
	}
	if namePrefix != nil {
		filters["name"] = *namePrefix
	}
	if minCount != nil {
		filters["min_count"] = *minCount
	}

	c.JSON(http.StatusOK, gin.H{
		"data": names,
		"meta": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+len(names) < total,
			"filters":  filters,
		},
	})
}

// Search handles GET /v1/names/search
func (h *NameHandler) Search(c *gin.Context) {
	// Parse required parameter
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "q parameter is required",
				Details: []FieldError{{Field: "q", Message: "required"}},
			},
		})
		return
	}

	if len(query) < 2 {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "search query must be at least 2 characters",
				Details: []FieldError{{Field: "q", Message: "must be at least 2 characters"}},
			},
		})
		return
	}

	// Parse optional parameters
	var country *string
	if c := c.Query("country"); c != "" {
		country = &c
	}

	var gender *string
	if g := c.Query("gender"); g != "" {
		if g != "M" && g != "F" {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: "gender must be M or F",
					Details: []FieldError{{Field: "gender", Message: "must be M or F"}},
				},
			})
			return
		}
		gender = &g
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit < 1 || limit > 1000 {
		limit = 100
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if offset < 0 {
		offset = 0
	}

	// Call service
	results, total, err := h.service.Search(
		c.Request.Context(),
		query,
		country,
		gender,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, service.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "not_found",
					Message: "Country not found",
				},
			})
			return
		}
		if errors.Is(err, service.ErrInvalidSearchQuery) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to search names",
			},
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data: results,
		Meta: PaginationMeta{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: offset+len(results) < total,
		},
	})
}

// GetByName handles GET /v1/names/:name
func (h *NameHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "name parameter is required",
			},
		})
		return
	}

	// Parse optional country filter
	var country *string
	if c := c.Query("country"); c != "" {
		country = &c
	}

	// Call service
	detail, err := h.service.GetByName(c.Request.Context(), name, country)
	if err != nil {
		if errors.Is(err, service.ErrNameNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "not_found",
					Message: "Name not found",
				},
			})
			return
		}
		if errors.Is(err, service.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "not_found",
					Message: "Country not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to get name details",
			},
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse{Data: detail})
}
