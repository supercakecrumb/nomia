package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/service"
)

// CountryHandler handles HTTP requests for country operations
type CountryHandler struct {
	service *service.CountryService
}

// NewCountryHandler creates a new country handler
func NewCountryHandler(service *service.CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

// FieldError represents a field validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// DataResponse represents a successful response with data
type DataResponse struct {
	Data interface{} `json:"data"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// Create handles POST /v1/countries
func (h *CountryHandler) Create(c *gin.Context) {
	var req model.CreateCountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "Invalid request body",
				Details: []FieldError{{Field: "body", Message: err.Error()}},
			},
		})
		return
	}

	country, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error: ErrorDetail{
					Code:    "validation_error",
					Message: err.Error(),
				},
			})
			return
		}
		if errors.Is(err, service.ErrCountryCodeExists) {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "conflict",
					Message: "Country code already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to create country",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, DataResponse{Data: country})
}

// List handles GET /v1/countries
func (h *CountryHandler) List(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	countries, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to list countries",
			},
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data: countries,
		Meta: PaginationMeta{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: offset+len(countries) < total,
		},
	})
}

// Get handles GET /v1/countries/:id
func (h *CountryHandler) Get(c *gin.Context) {
	idOrCode := c.Param("id")

	country, err := h.service.GetWithStats(c.Request.Context(), idOrCode)
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to get country",
			},
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse{Data: country})
}

// Update handles PATCH /v1/countries/:id
func (h *CountryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "Invalid country ID",
			},
		})
		return
	}

	var req model.UpdateCountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "Invalid request body",
				Details: []FieldError{{Field: "body", Message: err.Error()}},
			},
		})
		return
	}

	country, err := h.service.Update(c.Request.Context(), id, &req)
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
				Message: "Failed to update country",
			},
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse{Data: country})
}

// Delete handles DELETE /v1/countries/:id
func (h *CountryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "Invalid country ID",
			},
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
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
		// Check for foreign key constraint error
		if err.Error() == "cannot delete country with associated datasets" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "conflict",
					Message: "Cannot delete country with associated datasets",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "Failed to delete country",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}
