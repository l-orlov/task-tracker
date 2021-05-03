package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/l-orlov/task-tracker/internal/handler"
	"github.com/l-orlov/task-tracker/internal/repository"
	"github.com/l-orlov/task-tracker/internal/repository/postgres"
	"github.com/l-orlov/task-tracker/internal/server"
	"github.com/l-orlov/task-tracker/internal/service"
	"github.com/l-orlov/task-tracker/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := &config.Config{}
	var err error
	if err = config.ReadFromFileAndSetEnv(os.Getenv("CONFIG_PATH"), cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	lg, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	db, err := postgres.ConnectToDB(cfg.PostgresDB)
	if err != nil {
		lg.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			lg.Errorf("failed to close db: %v", err)
		}
	}()

	if cfg.PostgresDB.MigrationMode {
		if err = postgres.MigrateSchema(db.DB, cfg.PostgresDB); err != nil {
			log.Fatalf("failed to do migration: %v", err)
		}
	}

	repo, err := repository.NewRepository(cfg, lg, db)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	mailerLogEntry := logrus.NewEntry(lg).WithFields(logrus.Fields{"source": "mailerService"})
	mailer := service.NewMailerService(cfg.Mailer, mailerLogEntry)
	defer mailer.Close()

	svc, err := service.NewService(cfg, lg, repo, mailer)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	h := handler.New(cfg, lg, svc)

	srv := server.New(cfg.Port, h.InitRoutes())
	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	lg.Infof("service started on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	lg.Info("service shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		lg.Errorf("failed to shut down: %v", err)
	}
}
