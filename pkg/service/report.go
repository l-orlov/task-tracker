package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type ReportService struct {
	repo repository.Report
}

func NewReportService(repo repository.Report) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetAllProjectsWithTasksSubtasks(ctx context.Context) ([]models.ProjectWithTasksSubtasks, error) {
	projectWithTasksSubtasksDTO, err := s.repo.GetAllProjectsWithTasksSubtasks(ctx)
	if err != nil {
		return nil, err
	}

	_ = projectWithTasksSubtasksDTO

	return nil, nil
}
