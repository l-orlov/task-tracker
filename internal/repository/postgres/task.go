package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/models"
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

func (r *TaskPostgres) CreateTaskToProject(ctx context.Context, projectID int64, task models.TaskToCreate) (int64, error) {
	taskCreateQuery := fmt.Sprintf(`
INSERT INTO %s (title, description, assignee_id, importance_status_id, progress_status_id)
values ($1, $2, $3, $4, $5) RETURNING id`, tasksTable)
	setTaskToProjectQuery := fmt.Sprintf(`INSERT INTO %s (project_id, task_id) values ($1, $2)`, projectsTasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(dbCtx, taskCreateQuery, task.Title, task.Description, task.AssigneeID, task.ImportanceStatusID, task.ProgressStatusID)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var taskID int64
	if err := row.Scan(&taskID); err != nil {
		return 0, err
	}

	if _, err = tx.ExecContext(dbCtx, setTaskToProjectQuery, projectID, taskID); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return taskID, nil
}

func (r *TaskPostgres) GetTaskByID(ctx context.Context, id int64) (models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s WHERE id=$1`, tasksTable)
	var task models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &task, query, id)

	return task, err
}

func (r *TaskPostgres) UpdateTask(ctx context.Context, id int64, task models.TaskToUpdate) error {
	query := fmt.Sprintf(`
UPDATE %s SET title = $1, description = $2, assignee_id = $3,
importance_status_id = $4, progress_status_id = $5 WHERE id = $6`, tasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, task.Title, task.Description, task.AssigneeID,
		task.ImportanceStatusID, task.ProgressStatusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) GetAllTasksToProject(ctx context.Context, projectID int64) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT ts.id, ts.title, ts.description, ts.assignee_id, ts.importance_status_id, ts.progress_status_id
FROM %s as pts inner join %s as ts on pts.task_id = ts.id where pts.project_id = $1`, projectsTasksTable, tasksTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, projectID)

	return tasks, err
}

func (r *TaskPostgres) GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT id, title, description, assignee_id, importance_status_id, progress_status_id
FROM %s
WHERE (id = $1 OR $1 is null) AND (title ILIKE $2 OR $2 is null) AND (description ILIKE $3 OR $3 is null) AND
(assignee_id = $4 OR $4 is null) AND (importance_status_id = $5 OR $5 is null) AND (progress_status_id = $6 OR $6 is null)
ORDER BY id ASC`, tasksTable)

	if params.Title != nil {
		*params.Title = "%%" + *params.Title + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, params.ID, params.Title, params.Description,
		params.AssigneeID, params.ImportanceStatusID, params.ProgressStatusID)

	return tasks, err
}

func (r *TaskPostgres) GetAllTasksWithProjectID(ctx context.Context) ([]models.TaskWithProjectID, error) {
	query := fmt.Sprintf(`
SELECT pts.project_id, ts.id, ts.title, ts.description,
ts.assignee_id, ts.importance_status_id, ts.progress_status_id
FROM %s as pts inner join %s as ts on pts.task_id = ts.id`, projectsTasksTable, tasksTable)
	var tasks []models.TaskWithProjectID

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query)

	return tasks, err
}

func (r *TaskPostgres) DeleteTask(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, tasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
