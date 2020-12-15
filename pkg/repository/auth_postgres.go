package repository

import (
	"fmt"

	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r * AuthPostgres) CreateUser(user models.User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, email, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)
	
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
