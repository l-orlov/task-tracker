package models

type (
	ImportanceStatusToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
		Name      string `json:"name" binding:"required"`
	}
	ImportanceStatusToUpdate struct {
		ID int64 `json:"id" binding:"required"`
		ImportanceStatusToCreate
	}
	ImportanceStatus struct {
		ID        int64  `json:"id" db:"id"`
		ProjectID uint64 `json:"projectId" db:"project_id"`
		Name      string `json:"name" db:"name"`
	}
)
