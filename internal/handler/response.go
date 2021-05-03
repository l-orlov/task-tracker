package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ierrors "github.com/l-orlov/task-tracker/internal/errors"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, err error) {
	logEntry := h.getLogEntry(c)

	if customErr, ok := err.(*ierrors.Error); ok {
		handleCustomError(c, logEntry, customErr)
		return
	}

	handleDefaultError(c, logEntry, err, statusCode)
}

func handleCustomError(c *gin.Context, logEntry *logrus.Entry, err *ierrors.Error) {
	var statusCode int

	if err.Level == ierrors.Business {
		logEntry.Debug(err)
		statusCode = http.StatusBadRequest
	} else {
		logEntry.Error(err)
		statusCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(statusCode, &errorResponse{
		Message: err.Error(),
		Detail:  err.Detail,
	})
}

func handleDefaultError(c *gin.Context, logEntry *logrus.Entry, err error, statusCode int) {
	errResp := &errorResponse{
		Message: err.Error(),
	}
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		logEntry.Debug(err)
		errResp.Detail = ierrors.DetailBusiness
	} else {
		logEntry.Error(err)
		errResp.Detail = ierrors.DetailServer
	}

	c.AbortWithStatusJSON(statusCode, errResp)
}
