package repository

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (int64, error)
		GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error)
		GetUserByID(ctx context.Context, id int64) (models.UserToGet, error)
		UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error
		GetAllUsers(ctx context.Context) ([]models.UserToGet, error)
		DeleteUser(ctx context.Context, id int64) error
	}

	ImportanceStatus interface {
		Create(ctx context.Context, status models.StatusToCreate) (int64, error)
	}

	ProgressStatus interface {
		Create(ctx context.Context, status models.StatusToCreate) (int64, error)
	}

	Project interface {
		Create(ctx context.Context, project models.ProjectToCreate) (int64, error)
		GetAll(ctx context.Context) ([]models.Project, error)
	}

	Task interface {
	}

	Subtask interface {
	}

	Repository struct {
		User
		ImportanceStatus
		ProgressStatus
		Project
		Task
		Subtask
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:             NewUserPostgres(db),
		ImportanceStatus: NewImportanceStatusPostgres(db),
		ProgressStatus:   NewProgressStatusPostgres(db),
		Project:          NewProjectPostgres(db),
	}
}
