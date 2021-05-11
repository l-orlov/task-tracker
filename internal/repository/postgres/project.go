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

func (r *ProjectPostgres) CreateProject(
	ctx context.Context, project models.ProjectToCreate, owner uint64,
) (uint64, error) {
	createProjectQuery := fmt.Sprintf(`
INSERT INTO %s (name, description) values ($1, $2) RETURNING id`, projectTable)
	addProjectUserQuery := fmt.Sprintf(`
INSERT INTO %s (project_id, user_id, is_owner) values ($1, $2, 'TRUE')`, projectUserTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	row := r.db.QueryRowContext(dbCtx, createProjectQuery, &project.Name, &project.Description)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	if _, err = tx.ExecContext(dbCtx, addProjectUserQuery, id, owner); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProjectPostgres) GetProjectByID(ctx context.Context, id uint64) (*models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description FROM %s WHERE id = $1`, projectTable)
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

func (r *ProjectPostgres) UpdateProject(ctx context.Context, project models.Project) error {
	query := fmt.Sprintf(`
UPDATE %s SET name = :name, description = :description WHERE id = :id`, projectTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(dbCtx, query, &project)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description FROM %s ORDER BY id ASC`, projectTable)
	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query)

	return projects, err
}

func (r *ProjectPostgres) GetAllProjectsToUser(ctx context.Context, userID uint64) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT p.id, p.name, p.description
FROM %s AS p INNER JOIN %s AS pu ON p.id = pu.project_id
WHERE pu.user_id = $1 ORDER BY p.id ASC`, projectTable, projectUserTable)
	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query, &userID)

	return projects, err
}

func (r *ProjectPostgres) GetAllProjectsWithParameters(ctx context.Context, params models.ProjectParams) ([]models.Project, error) {
	query := fmt.Sprintf(`
SELECT id, name, description FROM %s
WHERE (id = $1 OR $1 is null) AND (name ILIKE $2 OR $2 is null) AND (description ILIKE $3 OR $3 is null)
ORDER BY id ASC`, projectTable)

	if params.Name != nil {
		*params.Name = "%%" + *params.Name + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var projects []models.Project

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projects, query, &params.ID, &params.Name, &params.Description)

	return projects, err
}

func (r *ProjectPostgres) DeleteProject(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, projectTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) AddUserToProject(ctx context.Context, projectID, userID uint64) error {
	query := fmt.Sprintf(`
INSERT INTO %s (project_id, user_id, is_owner) values ($1, $2, 'FALSE')`, projectUserTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &projectID, &userID); err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *ProjectPostgres) GetAllProjectUsers(ctx context.Context, projectID uint64) ([]models.ProjectUser, error) {
	query := fmt.Sprintf(`
SELECT u.id, u.email, u.firstname, u.lastname, pu.is_owner
FROM %s AS u INNER JOIN %s AS pu ON u.id = pu.user_id
WHERE pu.project_id = $1 ORDER BY u.id ASC`, userTable, projectUserTable)
	var users []models.ProjectUser

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.SelectContext(dbCtx, &users, query, &projectID); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ProjectPostgres) DeleteUserFromProject(ctx context.Context, projectID, userID uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE project_id = $1 AND user_id = $2`, projectUserTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &projectID, &userID); err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) GetProjectBoard(ctx context.Context, projectID uint64) (jsonData []byte, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s($1)`, fnGetProjectBoard)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err = r.db.QueryRowContext(dbCtx, query, &projectID).Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (r *ProjectPostgres) GetProgressStatusTask(ctx context.Context, taskID uint64) (jsonData []byte, err error) {
	query := fmt.Sprintf(`SELECT * FROM get_progress_status_task($1)`)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err = r.db.QueryRowContext(dbCtx, query, &taskID).Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
