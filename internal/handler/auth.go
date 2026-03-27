package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/contextx"
	"go-admin/server/internal/service"
	"go-admin/server/pkg/response"
)

type AuthHandler struct {
	auth *service.AuthService
	rbac *service.RBACService
}

func NewAuthHandler(auth *service.AuthService, rbac *service.RBACService) *AuthHandler {
	return &AuthHandler{auth: auth, rbac: rbac}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload service.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, tokens, err := h.auth.Login(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(contextx.UserIDKey, user.ID)
	c.Set(contextx.UsernameKey, user.Username)
	_, permissions, menus, err := h.rbac.GetProfile(c.Request.Context(), user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
		"permissions": permissions,
		"menus":       menus,
		"tokens":      tokens,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, tokens, err := h.auth.Refresh(c.Request.Context(), payload.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	_, permissions, menus, err := h.rbac.GetProfile(c.Request.Context(), user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"tokens": tokens, "permissions": permissions, "menus": menus})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get(contextx.UserIDKey)
	if err := h.auth.Logout(c.Request.Context(), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get(contextx.UserIDKey)
	user, permissions, menus, err := h.rbac.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"roles":    user.Roles,
		},
		"permissions": permissions,
		"menus":       menus,
	})
}
