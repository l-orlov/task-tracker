package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/internal/models"
	"golang.org/x/sync/errgroup"
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

func (h *Handler) UpdateProject(c *gin.Context) {
	var project models.ProjectToUpdate
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

func (h *Handler) GetAllProjectsWithTasks(c *gin.Context) {
	var (
		projects []models.Project
		tasks    []models.Task
	)

	g, gCtx := errgroup.WithContext(c)

	g.Go(func() error {
		var err error
		projects, err = h.svc.Project.GetAllProjects(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		tasks, err = h.svc.Task.GetAllTasks(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	tasksToProject := make(map[uint64][]models.Task)
	for _, task := range tasks {
		tasksToProject[task.ProjectID] = append(tasksToProject[task.ProjectID], task)
	}

	projectsWithTasks := make([]models.ProjectWithTasks, len(projects))
	for i := range projects {
		projectsWithTasks[i] = projects[i].ToProjectWithTasks()

		var tasks []models.Task
		var ok bool
		if tasks, ok = tasksToProject[projectsWithTasks[i].ID]; !ok {
			tasks = []models.Task{}
		}

		projectsWithTasks[i].Tasks = tasks
	}

	c.JSON(http.StatusOK, projectsWithTasks)
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
