package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type ProjectProgressStatusService struct {
	repo repository.ProjectProgressStatus
}

func NewProjectProgressStatusService(repo repository.ProjectProgressStatus) *ProjectProgressStatusService {
	return &ProjectProgressStatusService{repo: repo}
}

func (s *ProjectProgressStatusService) Add(ctx context.Context, projectID uint64, statusID int64) (int64, error) {
	return s.repo.Add(ctx, projectID, statusID)
}

func (s *ProjectProgressStatusService) GetByID(ctx context.Context, id int64) (*models.ProjectProgressStatus, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectProgressStatusService) GetAll(ctx context.Context) ([]models.ProjectProgressStatus, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectProgressStatusService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
