package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type ProgressStatusPostgres struct {
	db *sqlx.DB
}

func NewProgressStatusPostgres(db *sqlx.DB) *ProgressStatusPostgres {
	return &ProgressStatusPostgres{db: db}
}

func (r *ProgressStatusPostgres) Create(ctx context.Context, status models.StatusToCreate) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, status.Name)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProgressStatusPostgres) GetByID(ctx context.Context, id int64) (models.Status, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id=$1", progressStatusTable)
	var status models.Status

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &status, query, id)

	return status, err
}

func (r *ProgressStatusPostgres) Update(ctx context.Context, id int64, status models.StatusToCreate) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1 WHERE id = $2", progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, status.Name, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProgressStatusPostgres) GetAll(ctx context.Context) ([]models.Status, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s", progressStatusTable)
	var statuses []models.Status

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &statuses, query)

	return statuses, err
}

func (r *ProgressStatusPostgres) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", progressStatusTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
