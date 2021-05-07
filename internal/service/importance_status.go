package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type ImportanceStatusService struct {
	repo repository.ImportanceStatus
}

func NewImportanceStatusService(repo repository.ImportanceStatus) *ImportanceStatusService {
	return &ImportanceStatusService{repo: repo}
}

func (s *ImportanceStatusService) Create(ctx context.Context, status models.ImportanceStatusToCreate) (int64, error) {
	return s.repo.Create(ctx, status)
}

func (s *ImportanceStatusService) GetByID(ctx context.Context, id int64) (*models.ImportanceStatus, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ImportanceStatusService) Update(ctx context.Context, id int64, status models.ImportanceStatusToCreate) error {
	return s.repo.Update(ctx, id, status)
}

func (s *ImportanceStatusService) GetAll(ctx context.Context) ([]models.ImportanceStatus, error) {
	return s.repo.GetAll(ctx)
}

func (s *ImportanceStatusService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
