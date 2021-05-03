package models

import "time"

type (
	SprintToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
	}
	Sprint struct {
		ID        uint64     `json:"id" db:"id"`
		ProjectID uint64     `json:"projectId" db:"project_id"`
		CreatedAt time.Time  `json:"createdAt" db:"created_at"`
		ClosedAt  *time.Time `json:"closedAt" db:"closed_at"`
	}
	SprintParams struct {
		ID        *uint64    `json:"id"`
		ProjectID *uint64    `json:"projectId"`
		CreatedAt *time.Time `json:"createdAt"`
		ClosedAt  *time.Time `json:"closedAt"`
	}
)
