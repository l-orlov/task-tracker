package service

import (
	"context"

	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/l-orlov/task-tracker/internal/repository"
)

type ProjectService struct {
	repo repository.Project
}

func NewProjectService(repo repository.Project) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project models.ProjectToCreate, owner uint64) (uint64, error) {
	return s.repo.CreateProject(ctx, project, owner)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint64) (*models.Project, error) {
	return s.repo.GetProjectByID(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project models.ProjectToUpdate) error {
	return s.repo.UpdateProject(ctx, project)
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAllProjects(ctx)
}

func (s *ProjectService) GetAllProjectsToUser(ctx context.Context, userID uint64) ([]models.Project, error) {
	return s.repo.GetAllProjectsToUser(ctx, userID)
}

func (s *ProjectService) GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error) {
	return s.repo.GetAllProjectsWithParameters(ctx, params)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uint64) error {
	return s.repo.DeleteProject(ctx, id)
}

func (s *ProjectService) AddUserToProject(ctx context.Context, projectID, userID uint64) error {
	return s.repo.AddUserToProject(ctx, projectID, userID)
}

func (s *ProjectService) GetAllProjectUsers(ctx context.Context, projectID uint64) ([]models.ProjectUser, error) {
	return s.repo.GetAllProjectUsers(ctx, projectID)
}

func (s *ProjectService) DeleteUserFromProject(ctx context.Context, projectID, userID uint64) error {
	return s.repo.DeleteUserFromProject(ctx, projectID, userID)
}
