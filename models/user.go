package models

type (
	UserToCreate struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	UserToSignIn struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	UserToGet struct {
		ID        int64  `json:"id" db:"id"`
		FirstName string `json:"firstName" db:"first_name"`
		LastName  string `json:"lastName" db:"last_name"`
		Email     string `json:"email" db:"email"`
	}
)
