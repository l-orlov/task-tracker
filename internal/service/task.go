package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTaskToProject(ctx context.Context, task models.TaskToCreate) (uint64, error) {
	return s.repo.CreateTaskToProject(ctx, task)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uint64) (*models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, task models.Task) error {
	return s.repo.UpdateTask(ctx, task)
}

func (s *TaskService) GetAllTasksToProject(ctx context.Context, id uint64) ([]models.Task, error) {
	return s.repo.GetAllTasksToProject(ctx, id)
}

func (s *TaskService) GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error) {
	return s.repo.GetAllTasksWithParameters(ctx, params)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAllTasks(ctx)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uint64) error {
	return s.repo.DeleteTask(ctx, id)
}
