package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type ProgressStatusService struct {
	repo repository.ImportanceStatus
}

func NewProgressStatusService(repo repository.ProgressStatus) *ProgressStatusService {
	return &ProgressStatusService{repo: repo}
}

func (s *ProgressStatusService) Create(ctx context.Context, status models.StatusToCreate) (int64, error) {
	return s.repo.Create(ctx, status)
}
