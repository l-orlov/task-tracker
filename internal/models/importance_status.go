package models

type (
	ImportanceStatusToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
		Name      string `json:"name" binding:"required"`
	}
	ImportanceStatus struct {
		ID        int64  `json:"id" binding:"required" db:"id"`
		ProjectID uint64 `json:"projectId" binding:"required" db:"project_id"`
		Name      string `json:"name" binding:"required" db:"name"`
	}
)
