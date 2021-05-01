package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/pkg/errors"
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

func (r *ProjectPostgres) CreateProject(ctx context.Context, project models.ProjectToCreate, owner uint64) (uint64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (name, description, owner) values ($1, $2, $3) RETURNING id`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, project.Name, project.Description, owner)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProjectPostgres) GetProjectByID(ctx context.Context, id uint64) (*models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description, owner FROM %s WHERE id=$1`, projectsTable)
	var project models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &project, query, &id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &project, nil
}

func (r *ProjectPostgres) UpdateProject(ctx context.Context, project models.ProjectToUpdate) error {
	query := fmt.Sprintf(`
UPDATE %s SET name = $1, description = $2, owner = $3 WHERE id = $4`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, project.Name, project.Description, project.Owner, project.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description, owner FROM %s`, projectsTable)
	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query)

	return projects, err
}

func (r *ProjectPostgres) GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description, owner FROM %s
WHERE (id = $1 OR $1 is null) AND (name ILIKE $2 OR $2 is null) AND
(description ILIKE $3 OR $3 is null) AND (owner = $4 OR $4 is null) 
ORDER BY id ASC`, projectsTable)

	if params.Name != nil {
		*params.Name = "%%" + *params.Name + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(
		dbCtx, &projects, query, params.ID, params.Name, params.Description, params.Owner,
	)

	return projects, err
}

func (r *ProjectPostgres) DeleteProject(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, projectsTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
