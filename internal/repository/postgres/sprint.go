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

type SprintPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewSprintPostgres(db *sqlx.DB, dbTimeout time.Duration) *SprintPostgres {
	return &SprintPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *SprintPostgres) CreateSprintToProject(ctx context.Context, sprint models.SprintToCreate) (uint64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (project_id) values ($1) RETURNING id`, sprintTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, &sprint.ProjectID)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SprintPostgres) GetSprintByID(ctx context.Context, id uint64) (*models.Sprint, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, created_at, closed_at FROM %s WHERE id= $1`, sprintTable)
	var sprint models.Sprint

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &sprint, query, &id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &sprint, nil
}

func (r *SprintPostgres) GetAllSprintsToProject(ctx context.Context, projectID uint64) ([]models.Sprint, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, created_at, closed_at FROM %s WHERE project_id=$1`, sprintTable)
	var sprints []models.Sprint

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &sprints, query, &projectID)

	return sprints, err
}

func (r *SprintPostgres) GetAllSprintsWithParameters(ctx context.Context, params models.SprintParams) ([]models.Sprint, error) {
	query := fmt.Sprintf(`
SELECT id, project_id, created_at, closed_at
FROM %s
WHERE (id = $1 OR $1 is null) AND (project_id = $2 OR $2 is null) AND
(created_at = $3 OR $3 is null) AND (closed_at = $4 OR $4 is null)
ORDER BY id ASC`, sprintTable)

	var sprints []models.Sprint

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &sprints, query, &params.ID, &params.ProjectID,
		&params.CreatedAt, &params.ClosedAt)

	return sprints, err
}

func (r *SprintPostgres) CloseSprint(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`
UPDATE %s SET closed_at = $1 WHERE id = $2`, sprintTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	closedAt := time.Now().Format(time.RFC3339)

	_, err := r.db.ExecContext(dbCtx, query, &closedAt, &id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SprintPostgres) DeleteSprint(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, sprintTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}

func (r *SprintPostgres) AddTaskToSprint(ctx context.Context, sprintID, taskID uint64) error {
	query := fmt.Sprintf(`
INSERT INTO %s (sprint_id, task_id) values ($1, $2)`, sprintTaskTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &sprintID, &taskID); err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *SprintPostgres) GetAllSprintTasks(ctx context.Context, sprintID uint64) ([]models.Task, error) {
	query := fmt.Sprintf(`
SELECT t.id, t.project_id, t.title, t.description,
t.assignee_id, t.importance_status_id, t.progress_status_id
FROM %s AS t INNER JOIN %s AS st ON t.id = st.task_id
WHERE st.sprint_id = $1`, taskTable, sprintTaskTable)
	var tasks []models.Task

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.SelectContext(dbCtx, &tasks, query, &sprintID); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *SprintPostgres) DeleteTaskFromSprint(ctx context.Context, sprintID, taskID uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE sprint_id = $1 AND task_id = $2`, sprintTaskTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &sprintID, &taskID); err != nil {
		return err
	}

	return nil
}
