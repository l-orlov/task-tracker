package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/models"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) CreateProject(c *gin.Context) {
	var project models.ProjectToCreate
	if err := c.BindJSON(&project); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Project.CreateProject(c, project)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	project, err := h.services.Project.GetProjectByID(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var project models.ProjectToUpdate
	if err := c.BindJSON(&project); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Project.UpdateProject(c, id, project); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllProjects(c *gin.Context) {
	projects, err := h.services.Project.GetAllProjects(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetAllProjectsWithParameters(c *gin.Context) {
	var params models.ProjectParams
	if err := c.BindJSON(&params); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	projects, err := h.services.Project.GetAllProjectsWithParameters(c, params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetAllProjectsWithTasksSubtasks(c *gin.Context) {
	var (
		projects           []models.Project
		tasksWithProjectID []models.TaskWithProjectID
		subtasksWithTaskID []models.SubtaskWithTaskID
	)

	g, gCtx := errgroup.WithContext(c)

	g.Go(func() error {
		var err error
		projects, err = h.services.Project.GetAllProjects(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		tasksWithProjectID, err = h.services.Task.GetAllTasksWithProjectID(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		subtasksWithTaskID, err = h.services.Subtask.GetAllSubtasksWithTaskID(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	subtasksToTask := make(map[int64][]models.Subtask)
	for _, subtaskWithTaskID := range subtasksWithTaskID {
		subtasksToTask[subtaskWithTaskID.TaskID] = append(subtasksToTask[subtaskWithTaskID.TaskID], subtaskWithTaskID.ToSubtask())
	}

	tasksToProject := make(map[int64][]models.TaskWithSubtasks)
	for _, taskWithProjectID := range tasksWithProjectID {
		taskWithSubtasks := taskWithProjectID.ToTaskWithSubtasks()
		taskWithSubtasks.Subtasks = subtasksToTask[taskWithSubtasks.ID]
		tasksToProject[taskWithProjectID.ProjectID] = append(tasksToProject[taskWithProjectID.ProjectID], taskWithSubtasks)
	}

	projectsWithTsWithSs := make([]models.ProjectWithTasksWithSubtasks, len(projects))
	for i := range projects {
		projectsWithTsWithSs[i] = projects[i].ToProjectWithTasksWithSubtasks()
		projectsWithTsWithSs[i].Tasks = tasksToProject[projectsWithTsWithSs[i].ID]
	}

	c.JSON(http.StatusOK, projectsWithTsWithSs)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Project.DeleteProject(c, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
