package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
)

// DatasetHandler handles dataset-related requests
type DatasetHandler struct {
	datasetRepo *repository.DatasetRepository
}

// NewDatasetHandler creates a new DatasetHandler
func NewDatasetHandler(datasetRepo *repository.DatasetRepository) *DatasetHandler {
	return &DatasetHandler{
		datasetRepo: datasetRepo,
	}
}

// List handles GET /v1/datasets
func (h *DatasetHandler) List(c *gin.Context) {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate limits
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Parse filters
	filters := &model.DatasetFilters{}

	if countryIDStr := c.Query("country"); countryIDStr != "" {
		countryID, err := uuid.Parse(countryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "validation_error",
					"message": "Invalid country ID format",
				},
			})
			return
		}
		filters.CountryID = &countryID
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := model.DatasetStatus(statusStr)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "validation_error",
					"message": "Invalid status value",
				},
			})
			return
		}
		filters.Status = &status
	}

	// Get datasets
	datasets, total, err := h.datasetRepo.List(c.Request.Context(), filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to retrieve datasets",
			},
		})
		return
	}

	// Build response
	c.JSON(http.StatusOK, gin.H{
		"data": datasets,
		"meta": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+len(datasets) < total,
		},
	})
}

// Get handles GET /v1/datasets/:id
func (h *DatasetHandler) Get(c *gin.Context) {
	// Parse dataset ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "Invalid dataset ID format",
			},
		})
		return
	}

	// Get dataset with country information
	dataset, err := h.datasetRepo.GetByIDWithCountry(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "dataset not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "not_found",
					"message": "Dataset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to retrieve dataset",
			},
		})
		return
	}

	// Build response with processing time
	response := gin.H{
		"id":            dataset.ID,
		"country_id":    dataset.CountryID,
		"country_code":  dataset.CountryCode,
		"country_name":  dataset.CountryName,
		"filename":      dataset.Filename,
		"file_path":     dataset.FilePath,
		"file_size":     dataset.FileSize,
		"status":        dataset.Status,
		"row_count":     dataset.RowCount,
		"error_message": dataset.ErrorMessage,
		"uploaded_by":   dataset.UploadedBy,
		"uploaded_at":   dataset.UploadedAt,
		"processed_at":  dataset.ProcessedAt,
	}

	if processingTime := dataset.ProcessingTimeSeconds(); processingTime != nil {
		response["processing_time_seconds"] = *processingTime
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// Delete handles DELETE /v1/datasets/:id
func (h *DatasetHandler) Delete(c *gin.Context) {
	// Parse dataset ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "Invalid dataset ID format",
			},
		})
		return
	}

	// Delete dataset (soft delete)
	err = h.datasetRepo.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "dataset not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "not_found",
					"message": "Dataset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to delete dataset",
			},
		})
		return
	}

	// Return 204 No Content
	c.Status(http.StatusNoContent)
}
