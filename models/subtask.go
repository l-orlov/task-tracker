package models

type (
	SubtaskToCreate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	SubtaskToUpdate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		CreationDate       string `json:"creationDate" binding:"required"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	Subtask struct {
		ID                 int64  `json:"id" db:"id"`
		Title              string `json:"title" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"creationDate" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	SubtaskParams struct {
		ID                 *int64  `json:"id"`
		Title              *string `json:"title"`
		Description        *string `json:"description"`
		CreationDate       *string `json:"creationDate"`
		AssigneeID         *int64  `json:"assigneeId"`
		ImportanceStatusID *int64  `json:"importanceStatusId"`
		ProgressStatusID   *int64  `json:"progressStatusId"`
	}
)
