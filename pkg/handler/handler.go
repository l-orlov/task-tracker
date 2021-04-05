package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/l-orlov/task-tracker/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sigh-up", h.CreateUser)
		auth.POST("/sigh-in", h.SignIn)
	}

	api := router.Group("/api/v1", h.UserIdentity)
	{
		users := api.Group("/users")
		{
			users.POST("/", h.CreateUser)
			users.GET("/", h.GetAllUsers)
			users.GET("/by-id/:id", h.GetUserByID)
			users.GET("/by-email-password", h.GetUserByEmailPassword)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}

		importanceStatuses := api.Group("/importance")
		{
			importanceStatuses.POST("/", h.CreateImportanceStatus)
			importanceStatuses.GET("/", h.GetAllImportanceStatuses)
			importanceStatuses.GET("/:id", h.GetImportanceStatusByID)
			importanceStatuses.PUT("/:id", h.UpdateImportanceStatus)
			importanceStatuses.DELETE("/:id", h.DeleteImportanceStatus)
		}

		progressStatuses := api.Group("/progress")
		{
			progressStatuses.POST("/", h.CreateProgressStatus)
			progressStatuses.GET("/", h.GetAllProgressStatuses)
			progressStatuses.GET("/:id", h.GetProgressStatusByID)
			progressStatuses.PUT("/:id", h.UpdateProgressStatus)
			progressStatuses.DELETE("/:id", h.DeleteProgressStatus)
		}

		projects := api.Group("/projects")
		{
			projects.POST("/", h.CreateProject)
			projects.GET("/", h.GetAllProjects)
			projects.GET("/with-params", h.GetAllProjectsWithParameters)
			projects.GET("/with-tasks-subtasks", h.GetAllProjectsWithTasksSubtasks)
			projects.GET("/by-id/:id", h.GetProjectByID)
			projects.PUT("/:id", h.UpdateProject)
			projects.DELETE("/:id", h.DeleteProject)
		}

		tasks := api.Group("tasks")
		{
			tasks.POST("/", h.CreateTaskToProject)
			tasks.GET("/", h.GetAllTasksToProject)
			tasks.GET("/with-params", h.GetAllTasksWithParameters)
			tasks.GET("/by-id/:id", h.GetTaskByID)
			tasks.PUT("/:id", h.UpdateTask)
			tasks.DELETE("/:id", h.DeleteTask)
		}

		subtasks := api.Group("subtasks")
		{
			subtasks.POST("/", h.CreateSubtaskToTask)
			subtasks.GET("/", h.GetAllSubtasksToTask)
			subtasks.GET("/with-params", h.GetAllSubtasksWithParameters)
			subtasks.GET("/by-id/:id", h.GetSubtaskByID)
			subtasks.PUT("/:id", h.UpdateSubtask)
			subtasks.DELETE("/:id", h.DeleteSubtask)
		}
	}

	return router
}
