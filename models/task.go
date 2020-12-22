package models

type (
	ProjectToCreate struct {
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"-" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	Project struct {
		ID                 int64  `json:"id" binding:"required" db:"id"`
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"creationDate" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	TaskToCreate struct {
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"-" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	Task struct {
		ID                 int64  `json:"id" binding:"required" db:"id"`
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"creationDate" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	SubtaskToCreate struct {
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"-" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	Subtask struct {
		ID                 int64  `json:"id" binding:"required" db:"id"`
		Title              string `json:"title" binding:"required" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"creationDate" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}
)
