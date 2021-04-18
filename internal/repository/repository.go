package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/l-orlov/task-tracker/internal/models"
	cacheredis "github.com/l-orlov/task-tracker/internal/repository/cache-redis"
	userpostgres "github.com/l-orlov/task-tracker/internal/repository/user-postgres"
	"github.com/sirupsen/logrus"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		UpdateUserPassword(ctx context.Context, userID uint64, password string) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
	}
	SessionCache interface {
		PutSessionAndAccessToken(session models.Session, refreshToken string) error
		GetSession(refreshToken string) (*models.Session, error)
		DeleteSession(refreshToken string) error
		DeleteUserToSession(userID, refreshToken string) error
		GetAccessTokenData(accessTokenID string) (refreshToken string, err error)
		DeleteAccessToken(accessTokenID string) error
		AddUserBlocking(fingerprint string) (int64, error)
		GetUserBlocking(fingerprint string) (int, error)
		DeleteUserBlocking(fingerprint string) error
	}
	VerificationCache interface {
		PutEmailConfirmToken(userID uint64, token string) error
		GetEmailConfirmTokenData(token string) (userID uint64, err error)
		DeleteEmailConfirmToken(token string) error
		PutPasswordResetConfirmToken(userID uint64, token string) error
		GetPasswordResetConfirmTokenData(token string) (userID uint64, err error)
		DeletePasswordResetConfirmToken(token string) error
	}
	Repository struct {
		User
		SessionCache
		VerificationCache
	}
)

func NewRepository(
	cfg *config.Config, log *logrus.Logger, db *sqlx.DB,
) (*Repository, error) {
	userRepo := userpostgres.NewUserPostgres(db, cfg.PostgresDB.Timeout.Duration())

	cacheLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "cache-redis"})
	cacheOptions := cacheredis.Options{
		AccessTokenLifetime:               int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
		RefreshTokenLifetime:              int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
		UserBlockingLifetime:              int(cfg.UserBlocking.Lifetime.Duration().Seconds()),
		EmailConfirmTokenLifetime:         int(cfg.Verification.EmailConfirmTokenLifetime.Duration().Seconds()),
		PasswordResetConfirmTokenLifetime: int(cfg.Verification.PasswordResetConfirmTokenLifetime.Duration().Seconds()),
	}
	cache := cacheredis.New(cfg.Redis, cacheLogEntry, cacheOptions)

	return &Repository{
		User:              userRepo,
		SessionCache:      cache,
		VerificationCache: cache,
	}, nil
}
