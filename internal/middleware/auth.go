package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/config"
	"go-admin/server/internal/contextx"
	"go-admin/server/internal/service"
	"go-admin/server/pkg/jwtx"
	"go-admin/server/pkg/response"
)

func JWTAuth(cfg *config.Config, rbac *service.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}
		claims, err := jwtx.Parse(cfg.JWTSecret, token)
		if err != nil {
			response.Error(c, 401, "token invalid")
			c.Abort()
			return
		}
		_, permissions, _, err := rbac.GetProfile(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Error(c, 401, "user not found")
			c.Abort()
			return
		}
		c.Set(contextx.UserIDKey, claims.UserID)
		c.Set(contextx.UsernameKey, claims.Username)
		c.Set(contextx.PermissionsKey, permissions)
		c.Next()
	}
}

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, _ := c.Get(contextx.PermissionsKey)
		permissions, _ := raw.([]string)
		for _, permission := range permissions {
			if permission == code {
				c.Next()
				return
			}
		}
		response.Error(c, 403, "forbidden")
		c.Abort()
	}
}
