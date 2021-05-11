package models

type (
	ProgressStatusToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
		Name      string `json:"name" binding:"required"`
		OrderNum  int    `json:"orderNum"`
	}
	ProgressStatus struct {
		ID        int64  `json:"id" binding:"required" db:"id"`
		ProjectID uint64 `json:"projectId" binding:"required" db:"project_id"`
		Name      string `json:"name" binding:"required" db:"name"`
		OrderNum  int    `json:"orderNum" db:"order_num"`
	}
)
