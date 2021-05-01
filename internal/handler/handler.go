package handler

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/l-orlov/task-tracker/internal/service"
	"github.com/sirupsen/logrus"
)

type (
	Options struct {
		AccessTokenCookieMaxAge  int
		RefreshTokenCookieMaxAge int
		SecureCookie             *securecookie.SecureCookie
	}
	Handler struct {
		cfg        *config.Config
		log        *logrus.Logger
		options    Options
		svc        *service.Service
		httpServer *http.Server
	}
)

func New(
	cfg *config.Config, log *logrus.Logger, svc *service.Service,
) *Handler {
	c := &Handler{
		cfg: cfg,
		log: log,
		options: Options{
			AccessTokenCookieMaxAge:  int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
			RefreshTokenCookieMaxAge: int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
			SecureCookie:             securecookie.New(cfg.Cookie.HashKey, cfg.Cookie.BlockKey),
		},
		svc: svc,
	}

	return c
}

func (h *Handler) InitRoutes() http.Handler {
	router := gin.New()

	// for static files
	router.Use(static.Serve("/", static.LocalFile("./static", true)))

	router.Use(h.InitMiddleware)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.CreateUser)
		auth.POST("/sign-in", h.SignIn)
		router.POST("/reset-password", h.ResetPassword)
		auth.POST("/validate-access-token", h.ValidateAccessToken)
		auth.POST("/refresh-session", h.RefreshSession)
		auth.POST("/logout", h.Logout)
	}

	router.POST("/confirm-email", h.ConfirmEmail)
	router.POST("/confirm-reset-password", h.ConfirmPasswordReset)

	api := router.Group("/api/v1", h.UserAuthorizationMiddleware)
	{
		users := api.Group("/users")
		{
			users.POST("/", h.CreateUser)
			users.GET("/", h.GetAllUsers)
			users.GET("/with-params", h.GetAllUsersWithParameters)
			users.GET("/by-id/:id", h.GetUserByID)
			users.PUT("/", h.UpdateUser)
			users.PUT("/set-password", h.SetUserPassword)
			users.PUT("/change-password", h.ChangeUserPassword)
			users.DELETE("/by-id/:id", h.DeleteUser)
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
			projects.GET("/with-tasks", h.GetAllProjectsWithTasks)
			projects.GET("/:id", h.GetProjectByID)
			projects.PUT("/", h.UpdateProject)
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
	}

	return CORS(router)
}
