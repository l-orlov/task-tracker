package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createProject(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	_ = id
}

func (h *Handler) getAllProjects(c *gin.Context) {

}

func (h *Handler) getProjectByID(c *gin.Context) {

}

func (h *Handler) updateProject(c *gin.Context) {

}

func (h *Handler) deleteProject(c *gin.Context) {

}
