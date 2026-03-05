package handler

import (
	"errors"
	"net/http"
	"strconv"

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

	var req service.CreateBugRequest

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
		if errors.Is(err, errors.New("bug not found")) {
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

	c.JSON(http.StatusOK, gin.H{
		"bugs":   bugs,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateBugStatusRequest holds status update data
type UpdateBugStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateBugStatus updates a bug's status
// PATCH /bugs/:id/status
func (h *BugHandler) UpdateBugStatus(c *gin.Context) {
	bugID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bug id"})
		return
	}

	var req UpdateBugStatusRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bugService.UpdateBugStatus(c.Request.Context(), bugID, req.Status); err != nil {
		h.logger.Error("failed to update bug status", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bug status updated"})
}
