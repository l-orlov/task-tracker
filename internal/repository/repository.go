package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository/postgres"
	"github.com/l-orlov/task-tracker/internal/repository/redis"
	"github.com/sirupsen/logrus"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		UpdateUserPassword(ctx context.Context, userID uint64, password string) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		GetAllUsersWithParameters(ctx context.Context, params models.UserParams) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
	}
	Project interface {
		CreateProject(ctx context.Context, project models.ProjectToCreate, owner uint64) (uint64, error)
		GetProjectByID(ctx context.Context, id uint64) (*models.Project, error)
		UpdateProject(ctx context.Context, project models.ProjectToUpdate) error
		GetAllProjects(ctx context.Context) ([]models.Project, error)
		GetAllProjectsToUser(ctx context.Context, userID uint64) ([]models.Project, error)
		GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error)
		DeleteProject(ctx context.Context, id uint64) error
		AddUserToProject(ctx context.Context, projectID, userID uint64) error
		GetAllProjectUsers(ctx context.Context, projectID uint64) ([]models.ProjectUser, error)
		DeleteUserFromProject(ctx context.Context, projectID, userID uint64) error
	}
	ImportanceStatus interface {
		Create(ctx context.Context, status models.ImportanceStatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (*models.ImportanceStatus, error)
		Update(ctx context.Context, status models.ImportanceStatusToUpdate) error
		GetAll(ctx context.Context) ([]models.ImportanceStatus, error)
		GetAllToProject(ctx context.Context, projectID uint64) ([]models.ImportanceStatus, error)
		Delete(ctx context.Context, id int64) error
	}
	ProgressStatus interface {
		Create(ctx context.Context, status models.ProgressStatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (*models.ProgressStatus, error)
		Update(ctx context.Context, id int64, status models.ProgressStatusToCreate) error
		GetAll(ctx context.Context) ([]models.ProgressStatus, error)
		Delete(ctx context.Context, id int64) error
	}
	Task interface {
		CreateTaskToProject(ctx context.Context, task models.TaskToCreate) (uint64, error)
		GetTaskByID(ctx context.Context, id uint64) (*models.Task, error)
		UpdateTask(ctx context.Context, task models.TaskToUpdate) error
		GetAllTasksToProject(ctx context.Context, id uint64) ([]models.Task, error)
		GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error)
		GetAllTasks(ctx context.Context) ([]models.Task, error)
		DeleteTask(ctx context.Context, id uint64) error
	}
	SessionCache interface {
		PutSessionAndAccessToken(session models.Session, refreshToken string) error
		GetSession(refreshToken string) (*models.Session, error)
		DeleteSession(refreshToken string) error
		DeleteUserToSession(userID, refreshToken string) error
		GetAccessTokenData(accessTokenID string) (refreshToken string, err error)
		DeleteAccessToken(accessTokenID string) error
		AddUserBlocking(fingerprint string) (int64, error)
		GetUserBlocking(fingerprint string) (int, error)
		DeleteUserBlocking(fingerprint string) error
	}
	VerificationCache interface {
		PutEmailConfirmToken(userID uint64, token string) error
		GetEmailConfirmTokenData(token string) (userID uint64, err error)
		DeleteEmailConfirmToken(token string) error
		PutPasswordResetConfirmToken(userID uint64, token string) error
		GetPasswordResetConfirmTokenData(token string) (userID uint64, err error)
		DeletePasswordResetConfirmToken(token string) error
	}
	Repository struct {
		User
		Project
		ImportanceStatus
		ProgressStatus
		Task
		SessionCache
		VerificationCache
	}
)

func NewRepository(
	cfg *config.Config, log *logrus.Logger, db *sqlx.DB,
) (*Repository, error) {
	dbTimeout := cfg.PostgresDB.Timeout.Duration()

	cacheLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "cache-redis"})
	cacheOptions := redis.Options{
		AccessTokenLifetime:               int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
		RefreshTokenLifetime:              int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
		UserBlockingLifetime:              int(cfg.UserBlocking.Lifetime.Duration().Seconds()),
		EmailConfirmTokenLifetime:         int(cfg.Verification.EmailConfirmTokenLifetime.Duration().Seconds()),
		PasswordResetConfirmTokenLifetime: int(cfg.Verification.PasswordResetConfirmTokenLifetime.Duration().Seconds()),
	}
	cache := redis.New(cfg.Redis, cacheLogEntry, cacheOptions)

	return &Repository{
		User:              postgres.NewUserPostgres(db, dbTimeout),
		Project:           postgres.NewProjectPostgres(db, dbTimeout),
		ImportanceStatus:  postgres.NewImportanceStatusPostgres(db, dbTimeout),
		ProgressStatus:    postgres.NewProgressStatusPostgres(db, dbTimeout),
		Task:              postgres.NewTaskPostgres(db, dbTimeout),
		SessionCache:      cache,
		VerificationCache: cache,
	}, nil
}
