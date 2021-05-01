package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/task-tracker/internal/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ctxUserID   = "userID"
	ctxLogEntry = "log-entry"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) InitMiddleware(c *gin.Context) {
	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(h.log).WithField("request-id", requestID)

	logEntry.Infof("%s: %s", c.Request.Method, c.Request.RequestURI)

	c.Set(ctxLogEntry, logEntry)
	c.Next()
}

func (h *Handler) UserAuthorizationMiddleware(c *gin.Context) {
	err := h.validateTokenCookieAndRefreshIfNeeded(c)
	if err == nil {
		return
	}
	getLogEntry(c).Debug(err)

	if err = h.validateTokenHeader(c); err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Next()
}

// validateTokenCookieAndRefreshIfNeeded gets accessToken from cookie and validate it.
// on success it puts accessToken data to ctx and returns nil.
// else it tries to refresh session by refresh token from cookie:
// - on success puts accessToken data to ctx and returns nil
// - on failure returns error.
func (h *Handler) validateTokenCookieAndRefreshIfNeeded(c *gin.Context) error {
	accessToken, err := h.Cookie(c, accessTokenCookieName)
	if err != nil {
		getLogEntry(c).Debug(err)
		return h.refreshSessionByRefreshTokenCookie(c)
	}

	accessTokenClaims, err := h.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		if !strings.Contains(err.Error(), "token is expired by") &&
			!errors.Is(err, service.ErrNotActiveAccessToken) {
			return err
		}

		return h.refreshSessionByRefreshTokenCookie(c)
	}

	return setUserIDForContext(c, accessTokenClaims.Subject)
}

func (h *Handler) refreshSessionByRefreshTokenCookie(c *gin.Context) error {
	refreshToken, err := h.Cookie(c, refreshTokenCookieName)
	if err != nil {
		return err
	}

	newAccessToken, newRefreshToken, err := h.svc.UserAuthorization.RefreshSession(refreshToken)
	if err != nil {
		return err
	}

	accessTokenClaims, err := h.svc.UserAuthorization.GetAccessTokenClaims(newAccessToken)
	if err != nil {
		return err
	}

	h.setTokensCookies(c, newAccessToken, newRefreshToken)

	return setUserIDForContext(c, accessTokenClaims.Subject)
}

// validateTokenHeader gets accessToken from header and validate it.
// on success it puts accessToken data to ctx and returns nil. else it returns error.
func (h *Handler) validateTokenHeader(c *gin.Context) error {
	header := c.GetHeader("Authorization")
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return ErrNotValidAuthorizationHeader
	}

	accessToken := headerParts[1]
	accessTokenClaims, err := h.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		return err
	}

	return setUserIDForContext(c, accessTokenClaims.Subject)
}

func setHandlerNameToLogEntry(c *gin.Context, handlerName string) {
	logEntryValue, _ := c.Get(ctxLogEntry)

	logEntry := logEntryValue.(*logrus.Entry).WithField("method", handlerName)
	c.Set(ctxLogEntry, logEntry)
}

func getLogEntry(c *gin.Context) *logrus.Entry {
	logEntryValue, _ := c.Get(ctxLogEntry)

	return logEntryValue.(*logrus.Entry)
}

func setUserIDForContext(c *gin.Context, userIDStr string) error {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return err
	}

	c.Set(ctxUserID, userID)

	return nil
}
