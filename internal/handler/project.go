package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
)

func (h *Handler) CreateProject(c *gin.Context) {
	var project models.ProjectToCreate
	if err := c.BindJSON(&project); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	owner, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	id, err := h.svc.Project.CreateProject(c, project, owner)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	project, err := h.svc.Project.GetProjectByID(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if project == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) GetProjectByIDToUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	project, err := h.svc.Project.GetProjectByID(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if project == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	var project models.Project
	if err := c.BindJSON(&project); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Project.UpdateProject(c, project); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllProjects(c *gin.Context) {
	projects, err := h.svc.Project.GetAllProjects(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if projects == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetAllProjectsToUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	projects, err := h.svc.Project.GetAllProjectsToUser(c, userID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if projects == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetAllProjectsWithParameters(c *gin.Context) {
	var params models.ProjectParams
	if err := c.BindJSON(&params); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	projects, err := h.svc.Project.GetAllProjectsWithParameters(c, params)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if projects == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	if err := h.svc.Project.DeleteProject(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) AddUserToProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidUserIDQueryParam)
		return
	}

	if err := h.svc.Project.AddUserToProject(c, id, userID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllProjectUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	users, err := h.svc.Project.GetAllProjectUsers(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if users == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUserFromProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidIDParameter)
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidUserIDQueryParam)
		return
	}

	if err := h.svc.Project.DeleteUserFromProject(c, id, userID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
