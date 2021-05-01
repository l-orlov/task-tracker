package models

type (
	TaskToCreate struct {
		ProjectID          uint64 `json:"projectId" binding:"required"`
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		AssigneeID         uint64 `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	TaskToUpdate struct {
		ID uint64 `json:"id" binding:"required"`
		TaskToCreate
	}

	Task struct {
		ID                 uint64 `json:"id" db:"id"`
		ProjectID          uint64 `json:"projectId" db:"project_id"`
		Title              string `json:"title" db:"title"`
		Description        string `json:"description" db:"description"`
		AssigneeID         uint64 `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	TaskParams struct {
		ID                 *uint64 `json:"id"`
		ProjectID          *uint64 `json:"projectId"`
		Title              *string `json:"title"`
		Description        *string `json:"description"`
		AssigneeID         *uint64 `json:"assigneeId"`
		ImportanceStatusID *int64  `json:"importanceStatusId"`
		ProgressStatusID   *int64  `json:"progressStatusId"`
	}
)
