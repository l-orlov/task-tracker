package models

type (
	ProgressStatusToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
		Name      string `json:"name" binding:"required"`
		OrderNum  int    `json:"orderNum"`
	}
	ProgressStatusToUpdate struct {
		ID int64 `json:"id" binding:"required"`
		ProgressStatusToCreate
	}
	ProgressStatus struct {
		ID        int64  `json:"id" db:"id"`
		ProjectID uint64 `json:"projectId" db:"project_id"`
		Name      string `json:"name" db:"name"`
		OrderNum  int    `json:"orderNum" db:"order_num"`
	}
)
