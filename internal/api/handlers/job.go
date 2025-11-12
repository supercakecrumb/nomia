package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/supercakecrumb/affirm-name/internal/model"
	"github.com/supercakecrumb/affirm-name/internal/repository"
)

// JobHandler handles job-related requests
type JobHandler struct {
	jobRepo *repository.JobRepository
}

// NewJobHandler creates a new JobHandler
func NewJobHandler(jobRepo *repository.JobRepository) *JobHandler {
	return &JobHandler{
		jobRepo: jobRepo,
	}
}

// Get handles GET /v1/jobs/:id
func (h *JobHandler) Get(c *gin.Context) {
	// Parse job ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "validation_error",
				"message": "Invalid job ID format",
			},
		})
		return
	}

	// Get job
	job, err := h.jobRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "not_found",
					"message": "Job not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to retrieve job",
			},
		})
		return
	}

	// Build response with processing time
	response := gin.H{
		"id":           job.ID,
		"dataset_id":   job.DatasetID,
		"type":         job.Type,
		"status":       job.Status,
		"attempts":     job.Attempts,
		"max_attempts": job.MaxAttempts,
		"created_at":   job.CreatedAt,
		"started_at":   job.StartedAt,
		"completed_at": job.CompletedAt,
	}

	// Add optional fields
	if job.LastError != nil {
		response["last_error"] = *job.LastError
	}

	if job.Status == model.JobStatusCompleted && job.Payload != nil {
		response["result"] = job.Payload
	}

	if processingTime := job.ProcessingTimeSeconds(); processingTime != nil {
		response["processing_time_seconds"] = *processingTime
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// List handles GET /v1/jobs
func (h *JobHandler) List(c *gin.Context) {
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
	filters := &model.JobFilters{}

	if statusStr := c.Query("status"); statusStr != "" {
		status := model.JobStatus(statusStr)
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

	if typeStr := c.Query("type"); typeStr != "" {
		jobType := model.JobType(typeStr)
		if !jobType.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "validation_error",
					"message": "Invalid type value",
				},
			})
			return
		}
		filters.Type = &jobType
	}

	if datasetIDStr := c.Query("dataset_id"); datasetIDStr != "" {
		datasetID, err := uuid.Parse(datasetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "validation_error",
					"message": "Invalid dataset ID format",
				},
			})
			return
		}
		filters.DatasetID = &datasetID
	}

	// Get jobs
	jobs, total, err := h.jobRepo.List(c.Request.Context(), filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to retrieve jobs",
			},
		})
		return
	}

	// Build response - simplified job list without full payload
	jobList := make([]gin.H, len(jobs))
	for i, job := range jobs {
		jobList[i] = gin.H{
			"id":           job.ID,
			"dataset_id":   job.DatasetID,
			"type":         job.Type,
			"status":       job.Status,
			"attempts":     job.Attempts,
			"created_at":   job.CreatedAt,
			"completed_at": job.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": jobList,
		"meta": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+len(jobs) < total,
		},
	})
}
