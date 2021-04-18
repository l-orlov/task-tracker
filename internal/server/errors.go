package server

import (
	"github.com/pkg/errors"
)

var (
	ErrNotValidIDParameter = errors.New("not valid id parameter")
	ErrEmptyEmailParameter = errors.New("empty email parameter")
	ErrEmptyTokenParameter = errors.New("empty token parameter")
	ErrUserNotFound        = errors.New("user not found")
)
