package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTaskToProject(ctx context.Context, projectID int64, task models.TaskToCreate) (int64, error) {
	return s.repo.CreateTaskToProject(ctx, projectID, task)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, task models.TaskToUpdate) error {
	return s.repo.UpdateTask(ctx, id, task)
}

func (s *TaskService) GetAllTasksToProject(ctx context.Context, id int64) ([]models.Task, error) {
	return s.repo.GetAllTasksToProject(ctx, id)
}

func (s *TaskService) GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error) {
	return s.repo.GetAllTasksWithParameters(ctx, params)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return s.repo.DeleteTask(ctx, id)
}
