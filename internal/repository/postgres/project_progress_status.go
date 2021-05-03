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

type ProjectProgressStatusPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewProjectProgressStatusPostgres(db *sqlx.DB, dbTimeout time.Duration) *ProjectProgressStatusPostgres {
	return &ProjectProgressStatusPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ProjectProgressStatusPostgres) Add(ctx context.Context, projectID uint64, statusID int64) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (project_id, progress_status_id) values ($1, $2) RETURNING id`, projectProgressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, &projectID, &statusID)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProjectProgressStatusPostgres) GetByID(ctx context.Context, id int64) (*models.ProjectProgressStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, progress_status_id FROM %s WHERE id=$1`, projectProgressStatusTable)
	var status models.ProjectProgressStatus

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

func (r *ProjectProgressStatusPostgres) GetAll(ctx context.Context) ([]models.ProjectProgressStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, progress_status_id FROM %s`, projectProgressStatusTable)
	var statuses []models.ProjectProgressStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query)

	return statuses, err
}

func (r *ProjectProgressStatusPostgres) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, projectProgressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}
