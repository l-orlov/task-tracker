package repository

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type (
	Authorization interface {
		CreateUser(ctx context.Context, user models.User) (int64, error)
		GetUser(ctx context.Context, email, password string) (models.User, error)
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
		Authorization
		ImportanceStatus
		ProgressStatus
		Project
		Task
		Subtask
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:    NewAuthPostgres(db),
		ImportanceStatus: NewImportanceStatusPostgres(db),
		ProgressStatus:   NewProgressStatusPostgres(db),
		Project:          NewProjectPostgres(db),
	}
}
