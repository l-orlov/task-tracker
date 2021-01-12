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
