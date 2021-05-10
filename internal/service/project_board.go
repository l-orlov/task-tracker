package service

import (
	"context"

	ierrors "github.com/l-orlov/task-tracker/internal/errors"
	"github.com/l-orlov/task-tracker/internal/models"

	"github.com/l-orlov/task-tracker/internal/repository"
)

type ProjectBoardService struct {
	repo repository.ProjectBoard
}

func NewProjectBoardService(repo repository.ProjectBoard) *ProjectBoardService {
	return &ProjectBoardService{repo: repo}
}

func (s *ProjectBoardService) GetProjectBoardBytes(ctx context.Context, projectID uint64) (jsonData []byte, err error) {
	return s.repo.GetProjectBoardBytes(ctx, projectID)
}

func (s *ProjectBoardService) GetProjectBoard(ctx context.Context, projectID uint64) (*models.ProjectBoard, error) {
	return s.repo.GetProjectBoard(ctx, projectID)
}

func (s *ProjectBoardService) UpdateProjectBoardParts(ctx context.Context, board models.ProjectBoard) error {
	if len(board) != 2 {
		return ierrors.NewBusiness(ErrWrongProjectBoardPartsNum, "")
	}

	return s.repo.UpdateProjectBoardParts(ctx, board)
}

func (s *ProjectBoardService) UpdateProjectBoardProgressStatuses(ctx context.Context, statuses models.ProjectBoardProgressStatuses) error {
	return s.repo.UpdateProjectBoardProgressStatuses(ctx, statuses)
}

func (s *ProjectBoardService) UpdateProjectBoardProgressStatusTasks(ctx context.Context, tasks models.ProjectBoardProgressStatusTasks) error {
	return s.repo.UpdateProjectBoardProgressStatusTasks(ctx, tasks)
}
