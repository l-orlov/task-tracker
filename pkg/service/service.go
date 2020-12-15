package service

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int64, error)
}

type Project interface {

}

type Task interface {

}

type Subtask interface {

}

type Service struct {
	Authorization
	Project
	Task
	Subtask
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
