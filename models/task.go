package models

type (
	ProjectToCreate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	ProjectToUpdate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		CreationDate       string `json:"creationDate" binding:"required"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	Project struct {
		ID                 int64  `json:"id" db:"id"`
		Title              string `json:"title" db:"title"`
		Description        string `json:"description" db:"description"`
		CreationDate       string `json:"creationDate" db:"creation_date"`
		AssigneeID         int64  `json:"assigneeId" db:"assignee_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
		ProgressStatusID   int64  `json:"progressStatusId" db:"progress_status_id"`
	}

	TaskToCreate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	TaskToUpdate struct {
		Title              string `json:"title" binding:"required"`
		Description        string `json:"description"`
		CreationDate       string `json:"creationDate" binding:"required"`
		AssigneeID         int64  `json:"assigneeId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
		ProgressStatusID   int64  `json:"progressStatusId" binding:"required"`
	}

	Task struct {
		ID                 int64  `json:"id" db:"id"`
		Title              string `json:"title" db:"title"`
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
