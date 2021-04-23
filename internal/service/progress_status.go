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

func (s *ProgressStatusService) Create(ctx context.Context, status models.StatusToCreate) (int64, error) {
	return s.repo.Create(ctx, status)
}

func (s *ProgressStatusService) GetByID(ctx context.Context, id int64) (models.Status, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProgressStatusService) Update(ctx context.Context, id int64, status models.StatusToCreate) error {
	return s.repo.Update(ctx, id, status)
}

func (s *ProgressStatusService) GetAll(ctx context.Context) ([]models.Status, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProgressStatusService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
