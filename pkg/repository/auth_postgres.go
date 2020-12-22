package repository

import (
	"context"
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

func (r *AuthPostgres) CreateUser(user models.User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, email, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)

	dbCtx, cancel := context.WithTimeout(context.TODO(), dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (models.User, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	var user models.User

	dbCtx, cancel := context.WithTimeout(context.TODO(), dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &user, query, email, password)

	return user, err
}

//func (r * AuthPostgres) GetUser(email, password string) (models.User, error) {
//	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
//
//	dbCtx, cancel := context.WithTimeout(context.TODO(), dbTimeout)
//	defer cancel()
//
//	row := r.db.QueryRowContext(dbCtx, query, email, password)
//	if err := row.Err(); err != nil {
//		return models.User{}, err
//	}
//
//	var user models.User
//	if err := row.Scan(&user.ID); err != nil {
//		return models.User{}, err
//	}
//
//	return user, nil
//}
