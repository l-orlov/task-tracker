package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type ProjectImportanceStatusService struct {
	repo repository.ProjectImportanceStatus
}

func NewProjectImportanceStatusService(repo repository.ProjectImportanceStatus) *ProjectImportanceStatusService {
	return &ProjectImportanceStatusService{repo: repo}
}

func (s *ProjectImportanceStatusService) Add(ctx context.Context, projectID uint64, statusID int64) (int64, error) {
	return s.repo.Add(ctx, projectID, statusID)
}

func (s *ProjectImportanceStatusService) GetByID(ctx context.Context, id int64) (*models.ProjectImportanceStatus, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectImportanceStatusService) GetAll(ctx context.Context) ([]models.ProjectImportanceStatus, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectImportanceStatusService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
