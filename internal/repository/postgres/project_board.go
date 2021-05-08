package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/l-orlov/task-tracker/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProjectBoardPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewProjectBoardPostgres(db *sqlx.DB, dbTimeout time.Duration) *ProjectBoardPostgres {
	return &ProjectBoardPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ProjectBoardPostgres) GetProjectBoardBytes(ctx context.Context, projectID uint64) (jsonData []byte, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s($1)`, fnGetProjectBoard)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err = r.db.QueryRowContext(dbCtx, query, &projectID).Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (r *ProjectBoardPostgres) GetProjectBoard(ctx context.Context, projectID uint64) (*models.ProjectBoard, error) {
	query := fmt.Sprintf(`SELECT * FROM %s($1)`, fnGetProjectBoard)

	var board models.ProjectBoard

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.QueryRowContext(dbCtx, query, &projectID).Scan(&board)
	if err != nil {
		return nil, err
	}

	return &board, nil
}

func (r *ProjectBoardPostgres) UpdateProjectBoardParts(ctx context.Context, board models.ProjectBoard) error {
	query := fmt.Sprintf(`SELECT * FROM %s($1)`, fnUpdateProjectBoardParts)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, &board)
	if err != nil {
		return err
	}

	return nil
}
