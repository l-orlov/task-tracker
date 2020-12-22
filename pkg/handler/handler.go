package handler

import (
	"github.com/LevOrlov5404/task-tracker/pkg/service"
	"github.com/gin-gonic/gin"
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
		auth.POST("/sigh-up", h.signUp)
		auth.POST("/sigh-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		projects := api.Group("/projects")
		{
			projects.POST("/", h.createProject)
			projects.GET("/", h.getAllProjects)
			projects.GET("/:id", h.getProjectByID)
			projects.PUT("/:id", h.updateProject)
			projects.DELETE("/:id", h.deleteProject)

			tasks := projects.Group(":id/tasks")
			{
				tasks.POST("/", h.createTask)
				tasks.GET("/", h.getAllTasks)
				tasks.GET("/:id", h.getTaskByID)
				tasks.PUT("/:id", h.updateTask)
				tasks.DELETE("/:id", h.deleteTask)

				subtasks := tasks.Group(":id/subtasks")
				{
					subtasks.POST("/", h.createSubtask)
					subtasks.GET("/", h.getAllSubtasks)
					subtasks.GET("/:id", h.getSubtaskByID)
					subtasks.PUT("/:id", h.updateSubtask)
					subtasks.DELETE("/:id", h.deleteSubtask)
				}
			}
		}
	}

	return router
}
