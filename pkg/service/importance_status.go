package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type ImportanceStatusService struct {
	repo repository.ImportanceStatus
}

func NewImportanceStatusService(repo repository.ImportanceStatus) *ImportanceStatusService {
	return &ImportanceStatusService{repo: repo}
}

func (s *ImportanceStatusService) Create(ctx context.Context, status models.StatusToCreate) (int64, error) {
	return s.repo.Create(ctx, status)
}
