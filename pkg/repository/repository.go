package repository

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type (
	Authorization interface {
		CreateUser(user models.User) (int64, error)
		GetUser(email, password string) (models.User, error)
	}

	Project interface {
	}

	Task interface {
	}

	Subtask interface {
	}

	Repository struct {
		Authorization
		Project
		Task
		Subtask
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
