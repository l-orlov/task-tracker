package repository

import (
	"context"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (int64, error)
		GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error)
		GetUserByID(ctx context.Context, id int64) (models.UserToGet, error)
		UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error
		GetAllUsers(ctx context.Context) ([]models.UserToGet, error)
		DeleteUser(ctx context.Context, id int64) error
	}

	ImportanceStatus interface {
		Create(ctx context.Context, status models.StatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (models.Status, error)
		Update(ctx context.Context, id int64, status models.StatusToCreate) error
		GetAll(ctx context.Context) ([]models.Status, error)
		Delete(ctx context.Context, id int64) error
	}

	ProgressStatus interface {
		Create(ctx context.Context, status models.StatusToCreate) (int64, error)
		GetByID(ctx context.Context, id int64) (models.Status, error)
		Update(ctx context.Context, id int64, status models.StatusToCreate) error
		GetAll(ctx context.Context) ([]models.Status, error)
		Delete(ctx context.Context, id int64) error
	}

	Project interface {
		CreateProject(ctx context.Context, project models.ProjectToCreate) (int64, error)
		GetProjectByID(ctx context.Context, id int64) (models.Project, error)
		UpdateProject(ctx context.Context, id int64, project models.ProjectToUpdate) error
		GetAllProjects(ctx context.Context) ([]models.Project, error)
		GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error)
		DeleteProject(ctx context.Context, id int64) error
	}

	Task interface {
		CreateTaskToProject(ctx context.Context, projectID int64, task models.TaskToCreate) (int64, error)
		GetTaskByID(ctx context.Context, id int64) (models.Task, error)
		UpdateTask(ctx context.Context, id int64, task models.TaskToUpdate) error
		GetAllTasksToProject(ctx context.Context, id int64) ([]models.Task, error)
		GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error)
		GetAllTasksWithProjectID(ctx context.Context) ([]models.TaskWithProjectID, error)
		DeleteTask(ctx context.Context, id int64) error
	}

	Subtask interface {
		CreateSubtaskToTask(ctx context.Context, taskID int64, subtask models.SubtaskToCreate) (int64, error)
		GetSubtaskByID(ctx context.Context, id int64) (models.Subtask, error)
		UpdateSubtask(ctx context.Context, id int64, subtask models.SubtaskToUpdate) error
		GetAllSubtasksToTask(ctx context.Context, id int64) ([]models.Subtask, error)
		GetAllSubtasksWithParameters(ctx context.Context, params models.SubtaskParams) ([]models.Subtask, error)
		GetAllSubtasksWithTaskID(ctx context.Context) ([]models.SubtaskWithTaskID, error)
		DeleteSubtask(ctx context.Context, id int64) error
	}

	Repository struct {
		User
		ImportanceStatus
		ProgressStatus
		Project
		Task
		Subtask
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:             NewUserPostgres(db),
		ImportanceStatus: NewImportanceStatusPostgres(db),
		ProgressStatus:   NewProgressStatusPostgres(db),
		Project:          NewProjectPostgres(db),
		Task:             NewTaskPostgres(db),
		Subtask:          NewSubtaskPostgres(db),
	}
}
