package handler

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) CreateTaskToProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid projectId query param")
		return
	}

	var task models.TaskToCreate
	if err := c.BindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Task.CreateTaskToProject(c, projectID, task)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetTaskByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	task, err := h.services.Task.GetTaskByID(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var task models.TaskToUpdate
	if err := c.BindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Task.UpdateTask(c, id, task); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllTasksToProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid projectId query param")
		return
	}

	tasks, err := h.services.Task.GetAllTasksToProject(c, projectID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetAllTasksWithParameters(c *gin.Context) {
	var params models.TaskParams
	if err := c.BindJSON(&params); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tasks, err := h.services.Task.GetAllTasksWithParameters(c, params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Task.DeleteTask(c, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
