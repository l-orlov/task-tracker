package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/models"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateTaskToProject(ctx context.Context, projectID int64, task models.TaskToCreate) (int64, error) {
	taskCreateQuery := fmt.Sprintf("INSERT INTO %s (title, description, creation_date, assignee_id, importance_status_id, progress_status_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", tasksTable)
	setTaskToProjectQuery := fmt.Sprintf("INSERT INTO %s (project_id, task_id) values ($1, $2)", projectsTasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(dbCtx, taskCreateQuery, task.Title, task.Description, time.Now().Format(time.RFC3339), task.AssigneeID, task.ImportanceStatusID, task.ProgressStatusID)
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
	query := fmt.Sprintf("SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id FROM %s WHERE id=$1", tasksTable)
	var task models.Task

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &task, query, id)

	return task, err
}

func (r *TaskPostgres) UpdateTask(ctx context.Context, id int64, task models.TaskToUpdate) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2, creation_date = $3, assignee_id = $4, importance_status_id = $5, progress_status_id = $6 WHERE id = $7", tasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, task.Title, task.Description, task.CreationDate, task.AssigneeID, task.ImportanceStatusID, task.ProgressStatusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) GetAllTasksToProject(ctx context.Context, projectID int64) ([]models.Task, error) {
	query := fmt.Sprintf(`SELECT ts.id, ts.title, ts.description, ts.creation_date,
		ts.assignee_id, ts.importance_status_id, ts.progress_status_id
		FROM %s as pts inner join %s as ts on pts.task_id = ts.id where pts.project_id = $1`, projectsTasksTable, tasksTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, projectID)

	if tasks == nil {
		return []models.Task{}, nil
	}

	return tasks, err
}

func (r *TaskPostgres) GetAllTasksWithParameters(ctx context.Context, params models.TaskParams) ([]models.Task, error) {
	query := fmt.Sprintf(`SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id
		FROM %s
		WHERE (id = $1 OR $1 is null) AND (title ILIKE $2 OR $2 is null) AND (description ILIKE $3 OR $3 is null) AND (creation_date = $4 OR $4 is null)
		AND (assignee_id = $5 OR $5 is null) AND (importance_status_id = $6 OR $6 is null) AND (progress_status_id = $7 OR $7 is null)
		ORDER BY id ASC`, tasksTable)

	if params.Title != nil {
		*params.Title = "%%" + *params.Title + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, params.ID, params.Title, params.Description,
		params.CreationDate, params.AssigneeID, params.ImportanceStatusID, params.ProgressStatusID)

	if tasks == nil {
		return []models.Task{}, nil
	}

	return tasks, err
}

func (r *TaskPostgres) GetAllTasksWithProjectID(ctx context.Context) ([]models.TaskWithProjectID, error) {
	query := fmt.Sprintf(`SELECT pts.project_id, ts.id, ts.title, ts.description, ts.creation_date,
		ts.assignee_id, ts.importance_status_id, ts.progress_status_id
		FROM %s as pts inner join %s as ts on pts.task_id = ts.id`, projectsTasksTable, tasksTable)
	var tasks []models.TaskWithProjectID

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query)

	if tasks == nil {
		return []models.TaskWithProjectID{}, nil
	}

	return tasks, err
}

// ToDo: add deleting subtasks
func (r *TaskPostgres) DeleteTask(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
