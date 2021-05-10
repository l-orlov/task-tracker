package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/models"
	"github.com/pkg/errors"
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
INSERT INTO %s (email, firstname, lastname, password)
VALUES ($1, $2, $3, $4) RETURNING id`, userTable)
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
SELECT id, email, firstname, lastname, password, is_email_confirmed, avatar_url
FROM %s WHERE email=$1`, userTable)
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
SELECT id, email, firstname, lastname, password, is_email_confirmed, avatar_url
FROM %s WHERE id=$1`, userTable)
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
UPDATE %s SET firstname = $1, lastname = $2, avatar_url = $3 WHERE id = $4`, userTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, &user.FirstName, &user.LastName,
		&user.AvatarURL, &user.ID)
	if err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *UserPostgres) UpdateUserPassword(ctx context.Context, userID uint64, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password = $1 WHERE id = $2`, userTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &password, &userID); err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *UserPostgres) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, firstname, lastname, is_email_confirmed, avatar_url FROM %s ORDER BY id ASC`, userTable)
	var users []models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.SelectContext(dbCtx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserPostgres) GetAllUsersWithParameters(ctx context.Context, params models.UserParams) ([]models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, firstname, lastname, is_email_confirmed, avatar_url FROM %s
WHERE (id = $1 OR $1 is null) AND (email ILIKE $2 OR $2 is null) AND (firstname ILIKE $3 OR $3 is null) AND
(lastname = $4 OR $4 is null) AND (is_email_confirmed = $5 OR $5 is null)
ORDER BY id ASC`, userTable)

	if params.Email != nil {
		*params.Email = "%%" + *params.Email + "%%"
	}

	if params.FirstName != nil {
		*params.FirstName = "%%" + *params.FirstName + "%%"
	}

	if params.LastName != nil {
		*params.LastName = "%%" + *params.LastName + "%%"
	}

	var users []models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &users, query, params.ID, params.Email,
		params.FirstName, params.LastName, params.IsEmailConfirmed)

	return users, err
}

func (r *UserPostgres) DeleteUser(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, userTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) ConfirmEmail(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_email_confirmed = true WHERE id = $1`, userTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &id); err != nil {
		return getDBError(err)
	}

	return nil
}
