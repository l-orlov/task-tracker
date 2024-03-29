package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) CreateImportanceStatus(c *gin.Context) {
	var status models.ImportanceStatusToCreate
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.ImportanceStatus.Create(c, status)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetImportanceStatusByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	status, err := h.svc.ImportanceStatus.GetByID(c, id)
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

func (h *Handler) UpdateImportanceStatus(c *gin.Context) {
	var status models.ImportanceStatus
	if err := c.BindJSON(&status); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ImportanceStatus.Update(c, status); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllImportanceStatuses(c *gin.Context) {
	statuses, err := h.svc.ImportanceStatus.GetAll(c)
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

func (h *Handler) GetAllImportanceStatusesToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("projectId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidProjectIDQueryParam)
		return
	}

	statuses, err := h.svc.ImportanceStatus.GetAllToProject(c, projectID)
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

func (h *Handler) DeleteImportanceStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.ImportanceStatus.Delete(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
