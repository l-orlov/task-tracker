package handler

import (
	"github.com/pkg/errors"
)

var (
	ErrNotValidAuthorizationHeader = errors.New("not valid Authorization header")
	ErrNotValidIDParameter         = errors.New("not valid id parameter")
	ErrNotValidProjectIDParameter  = errors.New("invalid projectId query param")
	ErrEmptyEmailParameter         = errors.New("empty email parameter")
	ErrEmptyTokenParameter         = errors.New("empty token parameter")
	ErrUserNotFound                = errors.New("user not found")
)
