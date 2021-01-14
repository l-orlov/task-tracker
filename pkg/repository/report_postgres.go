package repository

import (
	"context"
	"fmt"
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type ReportPostgres struct {
	db *sqlx.DB
}

func NewReportPostgres(db *sqlx.DB) *ReportPostgres {
	return &ReportPostgres{db: db}
}

//3.	Получить все проекты с задачами и подзадачами

func (r *ReportPostgres) GetAllProjectsWithTasksSubtasks(ctx context.Context) ([]models.ProjectWithTasksSubtasksDTO, error) {
	query := fmt.Sprintf(`SELECT p.id, p.title, p.description, p.creation_date, p.assignee_id, p.importance_status_id, p.progress_status_id,
		t.id, t.title, t.description, t.creation_date, t.assignee_id, t.importance_status_id, t.progress_status_id,
		s.id, s.title, s.description, s.creation_date, s.assignee_id, s.importance_status_id, s.progress_status_id
		FROM %s as p
		LEFT JOIN %s as pts on p.id        = pts.project_id
		LEFT JOIN %s as t   on t.id        = pts.task_id
		LEFT JOIN %s as tss on tss.task_id = t.id
		LEFT JOIN %s as s   on s.id        = tss.subtask_id`, projectsTable, projectsTasksTable, tasksTable, tasksSubtasksTable, subtasksTable)
	var projectWithTasksSubtasks []models.ProjectWithTasksSubtasksDTO

	fmt.Println(query)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &projectWithTasksSubtasks, query)

	return projectWithTasksSubtasks, err
}
