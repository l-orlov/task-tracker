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

type ProjectImportanceStatusPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewProjectImportanceStatusPostgres(db *sqlx.DB, dbTimeout time.Duration) *ProjectImportanceStatusPostgres {
	return &ProjectImportanceStatusPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ProjectImportanceStatusPostgres) Add(ctx context.Context, projectID uint64, statusID int64) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (project_id, importance_status_id) values ($1, $2) RETURNING id`, projectImportanceStatusTable)

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

func (r *ProjectImportanceStatusPostgres) GetByID(ctx context.Context, id int64) (*models.ProjectImportanceStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, importance_status_id FROM %s WHERE id=$1`, projectImportanceStatusTable)
	var status models.ProjectImportanceStatus

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

func (r *ProjectImportanceStatusPostgres) GetAll(ctx context.Context) ([]models.ProjectImportanceStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, importance_status_id FROM %s`, projectImportanceStatusTable)
	var statuses []models.ProjectImportanceStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query)

	return statuses, err
}

func (r *ProjectImportanceStatusPostgres) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, projectImportanceStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}
