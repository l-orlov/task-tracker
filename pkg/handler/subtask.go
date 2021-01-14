package handler

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) CreateSubtaskToTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Query("taskId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid taskId query param")
		return
	}

	var subtask models.SubtaskToCreate
	if err := c.BindJSON(&subtask); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Subtask.CreateSubtaskToTask(c, taskID, subtask)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetSubtaskByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	subtask, err := h.services.Subtask.GetSubtaskByID(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subtask)
}

func (h *Handler) UpdateSubtask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var subtask models.SubtaskToUpdate
	if err := c.BindJSON(&subtask); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Subtask.UpdateSubtask(c, id, subtask); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllSubtasksToTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Query("taskId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid taskId query param")
		return
	}

	subtasks, err := h.services.Subtask.GetAllSubtasksToTask(c, taskID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subtasks)
}

func (h *Handler) GetAllSubtasksWithParameters(c *gin.Context) {
	var params models.SubtaskParams
	if err := c.BindJSON(&params); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subtasks, err := h.services.Subtask.GetAllSubtasksWithParameters(c, params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subtasks)
}

func (h *Handler) DeleteSubtask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Subtask.DeleteSubtask(c, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
