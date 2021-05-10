package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) GetProjectBoard(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("projectId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidProjectIDQueryParam)
		return
	}

	board, err := h.svc.ProjectBoard.GetProjectBoardBytes(c, projectID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if board == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.Data(200, "application/json", board)
}

func (h *Handler) UpdateProjectBoardParts(c *gin.Context) {
	var board models.ProjectBoard
	if err := c.BindJSON(&board); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ProjectBoard.UpdateProjectBoardParts(c, board); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateProjectBoardProgressStatuses(c *gin.Context) {
	var statuses models.ProjectBoardProgressStatuses
	if err := c.BindJSON(&statuses); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ProjectBoard.UpdateProjectBoardProgressStatuses(c, statuses); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateProjectBoardProgressStatusTasks(c *gin.Context) {
	var tasks models.ProjectBoardProgressStatusTasks
	if err := c.BindJSON(&tasks); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ProjectBoard.UpdateProjectBoardProgressStatusTasks(c, tasks); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
