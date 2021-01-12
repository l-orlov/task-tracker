package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type ImportanceStatusPostgres struct {
	db *sqlx.DB
}

func NewImportanceStatusPostgres(db *sqlx.DB) *ImportanceStatusPostgres {
	return &ImportanceStatusPostgres{db: db}
}

func (r *ImportanceStatusPostgres) Create(ctx context.Context, status models.StatusToCreate) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", importanceStatusTable)

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
