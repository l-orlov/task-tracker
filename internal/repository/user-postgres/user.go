package user_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	ierrors "github.com/l-orlov/task-tracker/internal/errors"
	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	usersTable = "users"
)

type UserPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewUserPostgres(db *sqlx.DB, dbTimeout time.Duration) *UserPostgres {
	return &UserPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (email, first_name, last_name, password)
VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)
	var err error

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query,
		&user.Email, &user.FirstName, &user.LastName, &user.Password)
	if err = row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id uint64
	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, first_name, last_name, password FROM %s WHERE email=$1`, usersTable)
	var user models.User
	var err error

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err = r.db.GetContext(dbCtx, &user, query, &email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserPostgres) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, first_name, last_name, password, is_email_confirmed FROM %s WHERE id=$1`, usersTable)
	var user models.User
	var err error

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err = r.db.GetContext(dbCtx, &user, query, &id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserPostgres) UpdateUser(ctx context.Context, user models.User) error {
	query := fmt.Sprintf(`
UPDATE %s SET first_name = $1, last_name = $2 WHERE id = $3`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, &user.FirstName, &user.LastName, &user.ID)
	if err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *UserPostgres) UpdateUserPassword(ctx context.Context, userID uint64, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password = $1 WHERE id = $2`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &password, &userID); err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *UserPostgres) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, first_name, last_name, is_email_confirmed FROM %s`, usersTable)
	var users []models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.SelectContext(dbCtx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserPostgres) DeleteUser(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) ConfirmEmail(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_email_confirmed = true WHERE id = $1`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return getDBError(err)
	}

	return nil
}

func getDBError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() < "50" { // business error
			return ierrors.NewBusiness(err, err.Detail)
		}

		return ierrors.New(err)
	}

	return err
}
