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

type ProgressStatusPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewProgressStatusPostgres(db *sqlx.DB, dbTimeout time.Duration) *ProgressStatusPostgres {
	return &ProgressStatusPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ProgressStatusPostgres) Create(ctx context.Context, status models.ProgressStatusToCreate) (int64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (project_id, name, order_num) values ($1, $2, $3) RETURNING id`, progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, &status.ProjectID, &status.Name, &status.OrderNum)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProgressStatusPostgres) GetByID(ctx context.Context, id int64) (*models.ProgressStatus, error) {
	query := fmt.Sprintf(`SELECT id, project_id, name, order_num FROM %s WHERE id=$1`, progressStatusTable)
	var status models.ProgressStatus

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

func (r *ProgressStatusPostgres) Update(ctx context.Context, status models.ProgressStatus) error {
	query := fmt.Sprintf(`
UPDATE %s SET name = :name, order_num = :order_num WHERE id = :id`, progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(dbCtx, query, &status)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProgressStatusPostgres) GetAll(ctx context.Context) ([]models.ProgressStatus, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, name, order_num FROM %s ORDER BY id ASC`, progressStatusTable)
	var statuses []models.ProgressStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query)

	return statuses, err
}

func (r *ProgressStatusPostgres) GetAllToProject(ctx context.Context, projectID uint64) ([]models.ProgressStatus, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, name, order_num FROM %s WHERE project_id = $1 ORDER BY id ASC`, progressStatusTable)
	var statuses []models.ProgressStatus

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query, &projectID)

	return statuses, err
}

func (r *ProgressStatusPostgres) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}
