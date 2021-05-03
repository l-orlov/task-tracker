package handler

import (
	"github.com/pkg/errors"
)

var (
	ErrNotValidAuthorizationHeader = errors.New("not valid Authorization header")
	ErrNotValidIDParameter         = errors.New("not valid id parameter")
	ErrNotValidProjectIDQueryParam = errors.New("not valid projectId query param")
	ErrNotValidUserIDQueryParam    = errors.New("not valid userId query param")
	ErrNotValidTaskIDQueryParam    = errors.New("not valid taskId query param")
	ErrEmptyEmailParameter         = errors.New("empty email parameter")
	ErrEmptyTokenParameter         = errors.New("empty token parameter")
	ErrUserNotFound                = errors.New("user not found")
)
