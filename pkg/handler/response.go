package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type customError struct {
	Message string `json:"message"`
}

func (err customError) Error() string {
	return err.Message
}

func newErrorRespnse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithError(statusCode, customError{
		Message: message,
	})
}
