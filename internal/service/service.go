package service

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
)

const (
	passwordAllowedLowerLetters = "abcdefghijklmnopqrstuvwxyz"
	passwordAllowedUpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	passwordAllowedDigits       = "0123456789"
)

type (
	RandomTokenGenerator interface {
		Generate(length, digitsNum, symbolsNum int, noUpper, allowRepeat bool) (string, error)
	}
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		SetUserPassword(ctx context.Context, userID uint64, password string) error
		ChangeUserPassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		GetAllUsersWithParameters(ctx context.Context, params models.UserParams) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
	}
	Project interface {
		CreateProject(ctx context.Context, project models.ProjectToCreate, owner uint64) (uint64, error)
		GetProjectByID(ctx context.Context, id uint64) (*models.Project, error)
		UpdateProject(ctx context.Context, project models.Project) error
		GetAllProjects(ctx context.Context) ([]models.Project, error)
		GetAllProjectsToUser(ctx context.Context, userID uint64) ([]models.Project, error)
		GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error)
		DeleteProject(ctx context.Context, id uint64) error
		AddUserToProject(ctx context.Context, projectID, userID uint64) error
		GetAllProjectUsers(ctx context.Context, projectID uint64) ([]models.ProjectUser, error)
		DeleteUserFromProject(ctx context.Context, projectID, userID uint64) error
	}
	ProjectBoard interface {
		GetProjectBoardBytes(ctx context.Context, projectID uint64) (jsonData []byte, err error)
		GetProjectBoard(ctx context.Context, projectID uint64) (*models.ProjectBoard, error)
		UpdateProjectBoardParts(ctx context.Context, board models.ProjectBoard) error
		UpdateProjectBoardProgressStatuses(ctx context.Context, statuses models.ProjectBoardProgressStatuses) error
		UpdateProjectBoardProgressStatusTasks(ctx context.Context, tasks models.ProjectBoardProgressStatusTasks) error
	}
	ImportanceStatus interface {
		Create(ctx context.Context, status models.ImportanceStatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (*models.ImportanceStatus, error)
		Update(ctx context.Context, status models.ImportanceStatus) error
		GetAll(ctx context.Context) ([]models.ImportanceStatus, error)
		GetAllToProject(ctx context.Context, projectID uint64) ([]models.ImportanceStatus, error)
		Delete(ctx context.Context, id int64) error
	}
	ProgressStatus interface {
		Create(ctx context.Context, status models.ProgressStatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (*models.ProgressStatus, error)
		Update(ctx context.Context, status models.ProgressStatus) error
		GetAll(ctx context.Context) ([]models.ProgressStatus, error)
		GetAllToProject(ctx context.Context, projectID uint64) ([]models.ProgressStatus, error)
		Delete(ctx context.Context, id int64) error
	}
	Task interface {
		CreateTaskToProject(ctx context.Context, task models.TaskToCreate) (uint64, error)
		GetTaskByID(ctx context.Context, id uint64) (*models.Task, error)
		UpdateTask(ctx context.Context, task models.Task) error
		GetAllTasksToProject(ctx context.Context, id uint64) ([]models.Task, error)
		GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error)
		GetAllTasks(ctx context.Context) ([]models.Task, error)
		DeleteTask(ctx context.Context, id uint64) error
	}
	UserAuthentication interface {
		AuthenticateUserByEmail(ctx context.Context, email, password, fingerprint string) (userID uint64, err error)
	}
	UserAuthorization interface {
		CreateSession(userID string) (accessToken, refreshToken string, err error)
		ValidateAccessToken(accessToken string) (*jwt.StandardClaims, error)
		RefreshSession(currentRefreshToken string) (accessToken, refreshToken string, err error)
		RevokeSession(accessToken string) error
		GetAccessTokenClaims(accessToken string) (*jwt.StandardClaims, error)
	}
	Verification interface {
		CreateEmailConfirmToken(userID uint64) (string, error)
		VerifyEmailConfirmToken(emailConfirmToken string) (userID uint64, err error)
		CreatePasswordResetConfirmToken(userID uint64) (string, error)
		VerifyPasswordResetConfirmToken(confirmToken string) (userID uint64, err error)
	}
	Mailer interface {
		SendEmailConfirm(toEmail, token string)
		SendResetPasswordConfirm(toEmail, token string)
	}
	Service struct {
		User
		Project
		ProjectBoard
		ImportanceStatus
		ProgressStatus
		Task
		UserAuthentication
		UserAuthorization
		Verification
		Mailer
	}
)

func NewService(
	cfg *config.Config, log *logrus.Logger,
	repo *repository.Repository, mailer Mailer,
) (*Service, error) {
	var generator RandomTokenGenerator
	var err error
	generator, err = password.NewGenerator(&password.GeneratorInput{
		LowerLetters: passwordAllowedLowerLetters,
		UpperLetters: passwordAllowedUpperLetters,
		Digits:       passwordAllowedDigits,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create random symbols generator")
	}

	authenticationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "authentication-svc"})
	verificationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "verification-svc"})

	return &Service{
		User:               NewUserService(repo.User, cfg.JWT.AccessTokenLifetime.Duration()),
		Project:            NewProjectService(repo.Project),
		ProjectBoard:       NewProjectBoardService(repo.ProjectBoard),
		ImportanceStatus:   NewImportanceStatusService(repo.ImportanceStatus),
		ProgressStatus:     NewProgressStatusService(repo.ProgressStatus),
		Task:               NewTaskService(repo.Task),
		UserAuthentication: NewAuthenticationService(cfg, authenticationLogEntry, repo),
		UserAuthorization:  NewAuthorizationService(cfg, repo),
		Verification:       NewVerificationService(verificationLogEntry, repo.VerificationCache, generator),
		Mailer:             mailer,
	}, nil
}
