package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) CreateTaskToProject(c *gin.Context) {
	var task models.TaskToCreate
	if err := c.BindJSON(&task); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.Task.CreateTaskToProject(c, task)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetTaskByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	task, err := h.svc.Task.GetTaskByID(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if task == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	var task models.TaskToUpdate
	if err := c.BindJSON(&task); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Task.UpdateTask(c, task); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllTasksToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("projectId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidProjectIDQueryParam)
		return
	}

	tasks, err := h.svc.Task.GetAllTasksToProject(c, projectID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if tasks == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetAllTasksWithParameters(c *gin.Context) {
	var params models.TaskParams
	if err := c.BindJSON(&params); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tasks, err := h.svc.Task.GetAllTasksWithParameters(c, params)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if tasks == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.Task.DeleteTask(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
