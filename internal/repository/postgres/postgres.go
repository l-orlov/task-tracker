package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/config"
)

const (
	userTable                    = "r_user"
	importanceStatusTable        = "s_importance_status"
	progressStatusTable          = "s_progress_status"
	projectTable                 = "r_project"
	projectUserTable             = "nn_project_user"
	projectImportanceStatusTable = "s_project_importance_status"
	projectProgressStatusTable   = "s_project_progress_status"
	taskTable                    = "r_task"
	sprintTable                  = "r_sprint"
	sprintTaskTable              = "nn_sprint_task"
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
