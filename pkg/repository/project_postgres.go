package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type ProjectPostgres struct {
	db *sqlx.DB
}

func NewProjectPostgres(db *sqlx.DB) *ProjectPostgres {
	return &ProjectPostgres{db: db}
}

func (r *ProjectPostgres) CreateProject(ctx context.Context, project models.ProjectToCreate) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, description, creation_date, assignee_id, importance_status_id, progress_status_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, project.Title, project.Description, time.Now(), project.AssigneeID, project.ImportanceStatusID, project.ProgressStatusID)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProjectPostgres) GetProjectByID(ctx context.Context, id int64) (models.Project, error) {
	query := fmt.Sprintf("SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id FROM %s WHERE id=$1", projectsTable)
	var project models.Project

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &project, query, id)

	return project, err
}

func (r *ProjectPostgres) UpdateProject(ctx context.Context, id int64, project models.ProjectToUpdate) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2, creation_date = $3, assignee_id = $4, importance_status_id = $5, progress_status_id = $6 WHERE id = $7", projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, project.Title, project.Description, project.CreationDate, project.AssigneeID, project.ImportanceStatusID, project.ProgressStatusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := fmt.Sprintf("SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id FROM %s", projectsTable)
	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query)

	return projects, err
}

// ToDo: add deleting tasks and subtasks
func (r *ProjectPostgres) DeleteProject(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
