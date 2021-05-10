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

type TaskPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewTaskPostgres(db *sqlx.DB, dbTimeout time.Duration) *TaskPostgres {
	return &TaskPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *TaskPostgres) CreateTaskToProject(ctx context.Context, task models.TaskToCreate) (uint64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (project_id, title, description, assignee_id, importance_status_id, progress_status_id)
values ($1, $2, $3, $4, $5, $6) RETURNING id`, taskTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, &task.ProjectID, &task.Title, &task.Description,
		&task.AssigneeID, &task.ImportanceStatusID, &task.ProgressStatusID)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TaskPostgres) GetTaskByID(ctx context.Context, id uint64) (*models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s WHERE id = $1`, taskTable)
	var task models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &task, query, &id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskPostgres) UpdateTask(ctx context.Context, task models.TaskToUpdate) error {
	query := fmt.Sprintf(`
UPDATE %s SET title = $1, description = $2, assignee_id = $3,
importance_status_id = $4, progress_status_id = $5 WHERE id = $6`, taskTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, &task.Title, &task.Description, &task.AssigneeID,
		&task.ImportanceStatusID, &task.ProgressStatusID, &task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) GetAllTasksToProject(ctx context.Context, projectID uint64) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s WHERE project_id=$1 ORDER BY id ASC`, taskTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, &projectID)

	return tasks, err
}

func (r *TaskPostgres) GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s
WHERE (id = $1 OR $1 is null) AND (project_id = $2 OR $2 is null) AND (title ILIKE $3 OR $3 is null) AND
(description ILIKE $4 OR $4 is null) AND (assignee_id = $5 OR $5 is null) AND 
(importance_status_id = $6 OR $6 is null) AND (progress_status_id = $7 OR $7 is null)
ORDER BY id ASC`, taskTable)

	if params.Title != nil {
		*params.Title = "%%" + *params.Title + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, &params.ID, &params.ProjectID, &params.Title,
		&params.Description, &params.AssigneeID, &params.ImportanceStatusID, &params.ProgressStatusID)

	return tasks, err
}

func (r *TaskPostgres) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s ORDER BY id ASC`, taskTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query)

	return tasks, err
}

func (r *TaskPostgres) DeleteTask(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, taskTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}
