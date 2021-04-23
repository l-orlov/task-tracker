package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/models"
)

type ProjectPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewProjectPostgres(db *sqlx.DB, dbTimeout time.Duration) *ProjectPostgres {
	return &ProjectPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ProjectPostgres) CreateProject(ctx context.Context, project models.ProjectToCreate) (int64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (title, description, assignee_id, importance_status_id, progress_status_id)
values ($1, $2, $3, $4, $5) RETURNING id`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, project.Title, project.Description,
		project.AssigneeID, project.ImportanceStatusID, project.ProgressStatusID)
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
	query := fmt.Sprintf(`
SELECT id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s WHERE id=$1`, projectsTable)
	var project models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &project, query, id)

	return project, err
}

func (r *ProjectPostgres) UpdateProject(ctx context.Context, id int64, project models.ProjectToUpdate) error {
	query := fmt.Sprintf(`
UPDATE %s SET title = $1, description = $2, assignee_id = $3, importance_status_id = $4,
progress_status_id = $5 WHERE id = $6`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, project.Title, project.Description,
		project.AssigneeID, project.ImportanceStatusID, project.ProgressStatusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s`, projectsTable)
	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query)

	return projects, err
}

func (r *ProjectPostgres) GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s
WHERE (id = $1 OR $1 is null) AND (title ILIKE $2 OR $2 is null) AND (description ILIKE $3 OR $3 is null) AND
(assignee_id = $4 OR $4 is null) AND (importance_status_id = $5 OR $5 is null) AND (progress_status_id = $6 OR $6 is null)
ORDER BY id ASC`, projectsTable)

	if params.Title != nil {
		*params.Title = "%%" + *params.Title + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query, params.ID, params.Title, params.Description,
		params.AssigneeID, params.ImportanceStatusID, params.ProgressStatusID)

	return projects, err
}

// ToDo: add deleting tasks
func (r *ProjectPostgres) DeleteProject(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
