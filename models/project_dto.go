package models

type (
	ProjectWithTasksSubtasksDTO struct {
		ProjectID                 int64  `db:"id"`
		ProjectTitle              string `db:"title"`
		ProjectDescription        string `db:"description"`
		ProjectCreationDate       string `db:"creation_date"`
		ProjectAssigneeID         int64  `db:"assignee_id"`
		ProjectImportanceStatusID int64  `db:"importance_status_id"`
		ProjectProgressStatusID   int64  `db:"progress_status_id"`
		TaskID                    int64  `db:"id"`
		TaskTitle                 string `db:"title"`
		TaskDescription           string `db:"description"`
		TaskCreationDate          string `db:"creation_date"`
		TaskAssigneeID            int64  `db:"assignee_id"`
		TaskImportanceStatusID    int64  `db:"importance_status_id"`
		TaskProgressStatusID      int64  `db:"progress_status_id"`
		SubtaskID                 int64  `db:"id"`
		SubtaskTitle              string `db:"title"`
		SubtaskDescription        string `db:"description"`
		SubtaskCreationDate       string `db:"creation_date"`
		SubtaskAssigneeID         int64  `db:"assignee_id"`
		SubtaskImportanceStatusID int64  `db:"importance_status_id"`
		SubtaskProgressStatusID   int64  `db:"progress_status_id"`
	}
)
