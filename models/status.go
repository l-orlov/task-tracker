package models

type (
	StatusToCreate struct {
		Name string `json:"name" binding:"required" db:"id"`
	}

	Status struct {
		ID   int64  `json:"id" db:"id"`
		Name string `json:"name" db:"name"`
	}
)
