package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) AddProjectImportanceStatus(c *gin.Context) {
	var status models.ProjectImportanceStatusToAdd
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.ProjectImportanceStatus.Add(c, status.ProjectID, status.ImportanceStatusID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetProjectImportanceStatusByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	status, err := h.svc.ProjectImportanceStatus.GetByID(c, id)
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

func (h *Handler) GetAllProjectImportanceStatuses(c *gin.Context) {
	statuses, err := h.svc.ProjectImportanceStatus.GetAll(c)
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

func (h *Handler) DeleteProjectImportanceStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.ProjectImportanceStatus.Delete(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) AddProjectProgressStatus(c *gin.Context) {
	var status models.ProjectProgressStatusToAdd
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.ProjectProgressStatus.Add(c, status.ProjectID, status.ProgressStatusID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetProjectProgressStatusByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	status, err := h.svc.ProjectProgressStatus.GetByID(c, id)
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

func (h *Handler) GetAllProjectProgressStatuses(c *gin.Context) {
	statuses, err := h.svc.ProjectProgressStatus.GetAll(c)
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

func (h *Handler) DeleteProjectProgressStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.ProjectProgressStatus.Delete(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
