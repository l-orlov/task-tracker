package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type SubtaskService struct {
	repo repository.Subtask
}

func NewSubtaskService(repo repository.Subtask) *SubtaskService {
	return &SubtaskService{repo: repo}
}

func (s *SubtaskService) CreateSubtaskToTask(ctx context.Context, taskID int64, subtask models.SubtaskToCreate) (int64, error) {
	return s.repo.CreateSubtaskToTask(ctx, taskID, subtask)
}

func (s *SubtaskService) GetSubtaskByID(ctx context.Context, id int64) (models.Subtask, error) {
	return s.repo.GetSubtaskByID(ctx, id)
}

func (s *SubtaskService) UpdateSubtask(ctx context.Context, id int64, subtask models.SubtaskToUpdate) error {
	return s.repo.UpdateSubtask(ctx, id, subtask)
}

func (s *SubtaskService) GetAllSubtasksToTask(ctx context.Context, id int64) ([]models.Subtask, error) {
	return s.repo.GetAllSubtasksToTask(ctx, id)
}

func (s *SubtaskService) GetAllSubtasksWithParameters(ctx context.Context, params models.SubtaskParams) ([]models.Subtask, error) {
	return s.repo.GetAllSubtasksWithParameters(ctx, params)
}

func (s *SubtaskService) GetAllSubtasksWithTaskID(ctx context.Context) ([]models.SubtaskWithTaskID, error) {
	return s.repo.GetAllSubtasksWithTaskID(ctx)
}

func (s *SubtaskService) DeleteSubtask(ctx context.Context, id int64) error {
	return s.repo.DeleteSubtask(ctx, id)
}
