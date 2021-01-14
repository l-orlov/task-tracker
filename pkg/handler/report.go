package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAllProjectsWithTasksSubtasks(c *gin.Context) {
	projectsWithTasksSubtasks, err := h.services.Report.GetAllProjectsWithTasksSubtasks(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, projectsWithTasksSubtasks)
}
