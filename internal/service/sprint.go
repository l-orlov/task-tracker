package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type SprintService struct {
	repo repository.Sprint
}

func NewSprintService(repo repository.Sprint) *SprintService {
	return &SprintService{repo: repo}
}

func (s *SprintService) CreateSprintToProject(ctx context.Context, sprint models.SprintToCreate) (uint64, error) {
	return s.repo.CreateSprintToProject(ctx, sprint)
}

func (s *SprintService) GetSprintByID(ctx context.Context, id uint64) (*models.Sprint, error) {
	return s.repo.GetSprintByID(ctx, id)
}

func (s *SprintService) GetAllSprintsToProject(ctx context.Context, projectID uint64) ([]models.Sprint, error) {
	return s.repo.GetAllSprintsToProject(ctx, projectID)
}

func (s *SprintService) GetAllSprintsWithParameters(ctx context.Context, params models.SprintParams) ([]models.Sprint, error) {
	return s.repo.GetAllSprintsWithParameters(ctx, params)
}

func (s *SprintService) CloseSprint(ctx context.Context, id uint64) error {
	return s.repo.CloseSprint(ctx, id)
}

func (s *SprintService) DeleteSprint(ctx context.Context, id uint64) error {
	return s.repo.DeleteSprint(ctx, id)
}

func (s *SprintService) AddTaskToSprint(ctx context.Context, sprintID, taskID uint64) error {
	return s.repo.AddTaskToSprint(ctx, sprintID, taskID)
}

func (s *SprintService) GetAllSprintTasks(ctx context.Context, sprintID uint64) ([]models.Task, error) {
	return s.repo.GetAllSprintTasks(ctx, sprintID)
}

func (s *SprintService) DeleteTaskFromSprint(ctx context.Context, sprintID, taskID uint64) error {
	return s.repo.DeleteTaskFromSprint(ctx, sprintID, taskID)
}
