package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/contextx"
	"go-admin/server/internal/model"
	"go-admin/server/internal/repository"
)

func RequestLogger(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := fmt.Sprintf("%d", start.UnixNano())
		c.Set(contextx.RequestIDKey, requestID)
		c.Next()

		if !shouldLogRequest(c) {
			return
		}

		userID, _ := c.Get(contextx.UserIDKey)
		username, _ := c.Get(contextx.UsernameKey)
		_ = repo.CreateLog(c.Request.Context(), &model.OperationLog{
			UserID:     toUint(userID),
			Username:   toString(username),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Action:     buildAction(c),
			LatencyMS:  time.Since(start).Milliseconds(),
			RequestID:  requestID,
		})
	}
}

func shouldLogRequest(c *gin.Context) bool {
	if c.Writer.Status() >= 400 {
		return false
	}

	switch c.FullPath() {
	case "/api/v1/auth/login", "/api/v1/auth/logout":
		return true
	case "/api/v1/auth/refresh", "/api/v1/auth/profile":
		return false
	}

	switch c.Request.Method {
	case "POST", "PUT", "PATCH", "DELETE":
		return len(c.FullPath()) >= len("/api/v1/system/") && c.FullPath()[:len("/api/v1/system/")] == "/api/v1/system/"
	default:
		return false
	}
}

func buildAction(c *gin.Context) string {
	switch c.FullPath() {
	case "/api/v1/auth/login":
		return "login"
	case "/api/v1/auth/logout":
		return "logout"
	}

	switch c.Request.Method {
	case "POST":
		return "create_or_save"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

func toUint(v any) uint {
	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

func toString(v any) string {
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}
