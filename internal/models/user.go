package models

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	UserToCreate struct {
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	UserToSignIn struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		Fingerprint string `json:"fingerprint" binding:"required"`
	}
	User struct {
		ID               uint64 `json:"id" binding:"required" db:"id"`
		Email            string `json:"email" binding:"required,email" db:"email"`
		FirstName        string `json:"firstName" binding:"required" db:"firstname"`
		LastName         string `json:"lastName" binding:"required" db:"lastname"`
		Password         string `json:"-" db:"password"`
		IsEmailConfirmed bool   `json:"isEmailConfirmed" db:"is_email_confirmed"`
		AvatarURL        string `json:"avatarURL" db:"avatar_url"`
	}
	UserPassword struct {
		ID       uint64 `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	UserPasswordToChange struct {
		ID          uint64 `json:"id" binding:"required"`
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	UserParams struct {
		ID               *uint64 `json:"id"`
		Email            *string `json:"email"`
		FirstName        *string `json:"firstName"`
		LastName         *string `json:"lastName"`
		IsEmailConfirmed *bool   `json:"isEmailConfirmed"`
	}
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
