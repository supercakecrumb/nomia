package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/service"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler creates a new UploadHandler
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// Upload handles POST /v1/datasets/upload
func (h *UploadHandler) Upload(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(service.MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "invalid_request",
				"message": "Failed to parse multipart form",
			},
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "File is required",
				"details": []gin.H{
					{
						"field":   "file",
						"message": "file field is required",
					},
				},
			},
		})
		return
	}

	// Get country_id from form
	countryIDStr := c.PostForm("country_id")
	if countryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "Country ID is required",
				"details": []gin.H{
					{
						"field":   "country_id",
						"message": "country_id field is required",
					},
				},
			},
		})
		return
	}

	// Parse country_id
	countryID, err := uuid.Parse(countryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "Invalid country ID format",
				"details": []gin.H{
					{
						"field":   "country_id",
						"message": "must be a valid UUID",
					},
				},
			},
		})
		return
	}

	// Validate file before processing
	if err := h.uploadService.ValidateFile(file); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if file.Size > service.MaxFileSize {
			statusCode = http.StatusRequestEntityTooLarge
		}

		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": err.Error(),
			},
		})
		return
	}

	// TODO: Get uploaded_by from authentication context
	// For now, use a placeholder
	uploadedBy := "admin@example.com"

	// Process upload
	result, err := h.uploadService.Upload(c.Request.Context(), &service.UploadRequest{
		File:       file,
		CountryID:  countryID,
		UploadedBy: uploadedBy,
	})
	if err != nil {
		// Check for specific error types
		if err.Error() == "country not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "not_found",
					"message": "Country not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to process upload",
			},
		})
		return
	}

	// Return 202 Accepted with dataset and job IDs
	c.JSON(http.StatusAccepted, gin.H{
		"data": result,
	})
}

// GetUploadLimits handles GET /v1/datasets/upload/limits
func (h *UploadHandler) GetUploadLimits(c *gin.Context) {
	limits := h.uploadService.GetUploadLimits()
	c.JSON(http.StatusOK, gin.H{
		"data": limits,
	})
}
