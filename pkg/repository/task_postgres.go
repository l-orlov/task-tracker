package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
	"time"
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

	row := tx.QueryRowContext(dbCtx, taskCreateQuery, task.Title, task.Description, time.Now(), task.AssigneeID, task.ImportanceStatusID, task.ProgressStatusID)
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
	query := fmt.Sprintf(`SELECT ts.id, ts.title, ts.description, ts.creation_date, ts.assignee_id, ts.importance_status_id, ts.progress_status_id
		FROM %s as pts inner join %s as ts on pts.task_id = ts.id where pts.project_id = $1`, projectsTasksTable, tasksTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &tasks, query, projectID)

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
