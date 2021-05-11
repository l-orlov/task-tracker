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

type ImportanceStatusPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewImportanceStatusPostgres(db *sqlx.DB, dbTimeout time.Duration) *ImportanceStatusPostgres {
	return &ImportanceStatusPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ImportanceStatusPostgres) Create(ctx context.Context, status models.ImportanceStatusToCreate) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (project_id, name) values ($1, $2) RETURNING id`, importanceStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, &status.ProjectID, &status.Name)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ImportanceStatusPostgres) GetByID(ctx context.Context, id int64) (*models.ImportanceStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, name FROM %s WHERE id=$1`, importanceStatusTable)
	var status models.ImportanceStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &status, query, &id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &status, nil
}

func (r *ImportanceStatusPostgres) Update(ctx context.Context, status models.ImportanceStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET name = :name WHERE id = :id`, importanceStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(dbCtx, query, &status)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImportanceStatusPostgres) GetAll(ctx context.Context) ([]models.ImportanceStatus, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, name FROM %s ORDER BY id ASC`, importanceStatusTable)
	var statuses []models.ImportanceStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query)

	return statuses, err
}

func (r *ImportanceStatusPostgres) GetAllToProject(ctx context.Context, projectID uint64) ([]models.ImportanceStatus, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, name FROM %s WHERE project_id = $1 ORDER BY id ASC`, importanceStatusTable)
	var statuses []models.ImportanceStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query, &projectID)

	return statuses, err
}

func (r *ImportanceStatusPostgres) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, importanceStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}
