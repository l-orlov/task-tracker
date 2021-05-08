package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type ProgressStatusService struct {
	repo repository.ProgressStatus
}

func NewProgressStatusService(repo repository.ProgressStatus) *ProgressStatusService {
	return &ProgressStatusService{repo: repo}
}

func (s *ProgressStatusService) Create(ctx context.Context, status models.ProgressStatusToCreate) (int64, error) {
	return s.repo.Create(ctx, status)
}

func (s *ProgressStatusService) GetByID(ctx context.Context, id int64) (*models.ProgressStatus, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProgressStatusService) Update(ctx context.Context, status models.ProgressStatusToUpdate) error {
	return s.repo.Update(ctx, status)
}

func (s *ProgressStatusService) GetAll(ctx context.Context) ([]models.ProgressStatus, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProgressStatusService) GetAllToProject(ctx context.Context, projectID uint64) ([]models.ProgressStatus, error) {
	return s.repo.GetAllToProject(ctx, projectID)
}

func (s *ProgressStatusService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
