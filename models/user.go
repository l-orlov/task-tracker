package models

type User struct {
	ID int64 `json:"-"`
	FirstName string `json:"firstName" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
