package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/LevOrlov5404/task-tracker/pkg/repository"
)

const (
	tokenTTL = 12 * time.Hour
)

type (
	AuthService struct {
		repo       repository.Authorization
		salt       string
		signingKey string
	}

	tokenClaims struct {
		jwt.StandardClaims
		UserID int64 `json:"userId"`
	}
)

func NewAuthService(repo repository.Authorization, salt, signingKey string) *AuthService {
	return &AuthService{
		repo:       repo,
		salt:       salt,
		signingKey: signingKey,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password, s.salt)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, email, generatePasswordHash(password, s.salt))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
