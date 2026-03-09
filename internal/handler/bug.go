package handler

import (
	"errors"
	"net/http"
	"strconv"

	"bug_triage/internal/dto"
	errortype "bug_triage/internal/err"
	"bug_triage/internal/middleware"
	"bug_triage/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BugHandler handles bug-related requests
type BugHandler struct {
	bugService *service.BugService
	logger     *zap.Logger
}

func NewBugHandler(bugService *service.BugService, logger *zap.Logger) *BugHandler {
	return &BugHandler{
		bugService: bugService,
		logger:     logger,
	}
}

// CreateBug handles bug creation
// POST /bugs
func (h *BugHandler) CreateBug(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req dto.CreateBugRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bug, err := h.bugService.CreateBug(c.Request.Context(), &req, userID)
	if err != nil {
		h.logger.Error("failed to create bug", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create bug"})
		return
	}

	c.JSON(http.StatusCreated, bug)
}

// GetBug retrieves a single bug
// GET /bugs/:id
func (h *BugHandler) GetBug(c *gin.Context) {
	bugID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bug id"})
		return
	}

	bug, err := h.bugService.GetBug(c.Request.Context(), bugID)
	if err != nil {
		if errors.Is(err, errortype.ErrBugNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bug not found"})
			return
		}
		h.logger.Error("failed to get bug", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get bug"})
		return
	}

	c.JSON(http.StatusOK, bug)
}

// ListBugs retrieves a paginated list of bugs
// GET /bugs?limit=20&offset=0
func (h *BugHandler) ListBugs(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	bugs, err := h.bugService.ListBugs(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list bugs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bugs"})
		return
	}

	// Convert models.Bug to dto.BugResponse
	bugResponses := make([]dto.BugResponse, len(bugs))
	for i, bug := range bugs {
		bugResponses[i] = dto.BugResponse{
			ID:          bug.ID,
			Title:       bug.Title,
			Description: bug.Description,
			Status:      bug.Status,
			Priority:    bug.Priority,
			Category:    bug.Category,
			ReporterID:  bug.ReporterID,
			CreatedAt:   bug.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	response := dto.BugsListResponse{
		Bugs:   bugResponses,
		Limit:  limit,
		Offset: offset,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateBugStatus updates a bug's status
// PATCH /bugs/:id/status
func (h *BugHandler) UpdateBugStatus(c *gin.Context) {
	bugID, err := strconv.ParseInt(c.Param("id"), 10, 64) // string, base (10 as base) , bit(64 bit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bug id"})
		return
	}

	var req dto.UpdateBugStatusRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bugService.UpdateBugStatus(c.Request.Context(), bugID, req.Status)
	if err != nil {
		if errors.Is(err, errortype.ErrBugNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bug not found"})
			return
		}
		h.logger.Error("failed to update bug status", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bug status updated"})
}
