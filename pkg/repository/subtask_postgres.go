package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type SubtaskPostgres struct {
	db *sqlx.DB
}

func NewSubtaskPostgres(db *sqlx.DB) *SubtaskPostgres {
	return &SubtaskPostgres{db: db}
}

func (r *SubtaskPostgres) CreateSubtaskToTask(ctx context.Context, taskID int64, subtask models.SubtaskToCreate) (int64, error) {
	subtaskCreateQuery := fmt.Sprintf("INSERT INTO %s (title, description, creation_date, assignee_id, importance_status_id, progress_status_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", subtasksTable)
	setSubtaskToProjectQuery := fmt.Sprintf("INSERT INTO %s (task_id, subtask_id) values ($1, $2)", tasksSubtasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(dbCtx, subtaskCreateQuery, subtask.Title, subtask.Description, time.Now().Format(time.RFC3339), subtask.AssigneeID, subtask.ImportanceStatusID, subtask.ProgressStatusID)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var subtaskID int64
	if err := row.Scan(&subtaskID); err != nil {
		return 0, err
	}

	if _, err = tx.ExecContext(dbCtx, setSubtaskToProjectQuery, taskID, subtaskID); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return subtaskID, nil
}

func (r *SubtaskPostgres) GetSubtaskByID(ctx context.Context, id int64) (models.Subtask, error) {
	query := fmt.Sprintf("SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id FROM %s WHERE id=$1", subtasksTable)
	var subtask models.Subtask

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &subtask, query, id)

	return subtask, err
}

func (r *SubtaskPostgres) UpdateSubtask(ctx context.Context, id int64, subtask models.SubtaskToUpdate) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2, creation_date = $3, assignee_id = $4, importance_status_id = $5, progress_status_id = $6 WHERE id = $7", subtasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, subtask.Title, subtask.Description, subtask.CreationDate, subtask.AssigneeID, subtask.ImportanceStatusID, subtask.ProgressStatusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubtaskPostgres) GetAllSubtasksToTask(ctx context.Context, projectID int64) ([]models.Subtask, error) {
	query := fmt.Sprintf(`SELECT ss.id, ss.title, ss.description, ss.creation_date,
		ss.assignee_id, ss.importance_status_id, ss.progress_status_id
		FROM %s as tss inner join %s as ss on tss.subtask_id = ss.id where tss.task_id = $1`, tasksSubtasksTable, subtasksTable)
	var subtasks []models.Subtask

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &subtasks, query, projectID)

	if subtasks == nil {
		return []models.Subtask{}, nil
	}

	return subtasks, err
}

func (r *SubtaskPostgres) GetAllSubtasksWithParameters(ctx context.Context, params models.SubtaskParams) ([]models.Subtask, error) {
	query := fmt.Sprintf(`SELECT id, title, description, creation_date, assignee_id, importance_status_id, progress_status_id
		FROM %s
		WHERE (id = $1 OR $1 is null) AND (title ILIKE $2 OR $2 is null) AND (description ILIKE $3 OR $3 is null) AND (creation_date = $4 OR $4 is null)
		AND (assignee_id = $5 OR $5 is null) AND (importance_status_id = $6 OR $6 is null) AND (progress_status_id = $7 OR $7 is null)
		ORDER BY id ASC`, subtasksTable)

	if params.Title != nil {
		*params.Title = "%%" + *params.Title + "%%"
	}

	if params.Description != nil {
		*params.Description = "%%" + *params.Description + "%%"
	}

	var subtasks []models.Subtask

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &subtasks, query, params.ID, params.Title, params.Description,
		params.CreationDate, params.AssigneeID, params.ImportanceStatusID, params.ProgressStatusID)

	if subtasks == nil {
		return []models.Subtask{}, nil
	}

	return subtasks, err
}

func (r *SubtaskPostgres) GetAllSubtasksWithTaskID(ctx context.Context) ([]models.SubtaskWithTaskID, error) {
	query := fmt.Sprintf(`SELECT tss.task_id, ss.id, ss.title, ss.description, ss.creation_date,
		ss.assignee_id, ss.importance_status_id, ss.progress_status_id
		FROM %s as tss inner join %s as ss on tss.subtask_id = ss.id`, tasksSubtasksTable, subtasksTable)
	var subtasks []models.SubtaskWithTaskID

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &subtasks, query)

	if subtasks == nil {
		return []models.SubtaskWithTaskID{}, nil
	}

	return subtasks, err
}

func (r *SubtaskPostgres) DeleteSubtask(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", subtasksTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
