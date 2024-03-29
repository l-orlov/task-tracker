package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) CreateProgressStatus(c *gin.Context) {
	var status models.ProgressStatusToCreate
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.ProgressStatus.Create(c, status)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetProgressStatusByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	status, err := h.svc.ProgressStatus.GetByID(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if status == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *Handler) UpdateProgressStatus(c *gin.Context) {
	var status models.ProgressStatus
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ProgressStatus.Update(c, status); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllProgressStatuses(c *gin.Context) {
	statuses, err := h.svc.ProgressStatus.GetAll(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if statuses == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (h *Handler) GetAllProgressStatusesToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("projectId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidProjectIDQueryParam)
		return
	}

	statuses, err := h.svc.ProgressStatus.GetAllToProject(c, projectID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if statuses == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (h *Handler) DeleteProgressStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.ProgressStatus.Delete(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
