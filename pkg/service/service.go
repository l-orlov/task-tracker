package service

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type (
	Authorization interface {
		CreateUser(user models.User) (int64, error)
		GenerateToken(email, password string) (string, error)
		ParseToken(token string) (int64, error)
	}

	Project interface {
	}

	Task interface {
	}

	Subtask interface {
	}

	Service struct {
		Authorization
		Project
		Task
		Subtask
	}
)

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
