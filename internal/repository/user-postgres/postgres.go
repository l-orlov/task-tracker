package user_postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/config"
)

func ConnectToDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.PostgresDB.URL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
