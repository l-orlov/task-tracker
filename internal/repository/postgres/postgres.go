package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/config"
	ierrors "github.com/l-orlov/task-tracker/internal/errors"
	"github.com/lib/pq"
)

const (
	userTable             = "r_user"
	projectTable          = "r_project"
	importanceStatusTable = "s_project_importance_status"
	progressStatusTable   = "s_project_progress_status"
	projectUserTable      = "nn_project_user"
	taskTable             = "r_task"
)

func ConnectToDB(cfg config.PostgresDB) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", initConnectionString(cfg))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initConnectionString(cfg config.PostgresDB) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Address.Host, cfg.Address.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
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
