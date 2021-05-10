package models

type (
	ProjectToCreate struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	ProjectToUpdate struct {
		ID uint64 `json:"id" binding:"required"`
		ProjectToCreate
	}
	Project struct {
		ID          uint64 `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
	}
	ProjectParams struct {
		ID          *uint64 `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	ProjectUser struct {
		ID        uint64 `json:"id" db:"id"`
		Email     string `json:"email" db:"email"`
		FirstName string `json:"firstName" db:"firstname"`
		LastName  string `json:"lastName" db:"lastname"`
		IsOwner   bool   `json:"isOwner" db:"is_owner"`
	}
)
