package service

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

type ProjectService struct {
	repo repository.Project
}

func NewProjectService(repo repository.Project) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project models.ProjectToCreate) (int64, error) {
	return s.repo.CreateProject(ctx, project)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id int64) (models.Project, error) {
	return s.repo.GetProjectByID(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id int64, project models.ProjectToUpdate) error {
	return s.repo.UpdateProject(ctx, id, project)
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAllProjects(ctx)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	return s.repo.DeleteProject(ctx, id)
}
