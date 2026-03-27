package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/config"
	"go-admin/server/internal/contextx"
	"go-admin/server/internal/model"
	"go-admin/server/internal/service"
	"go-admin/server/internal/utils"
	"go-admin/server/pkg/response"
)

type SystemHandler struct {
	system *service.SystemService
	cfg    *config.Config
}

func NewSystemHandler(system *service.SystemService, cfg *config.Config) *SystemHandler {
	return &SystemHandler{system: system, cfg: cfg}
}

func (h *SystemHandler) ListUsers(c *gin.Context) {
	var query service.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	users, err := h.system.ListUsers(c.Request.Context(), operatorUserID.(uint), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

func (h *SystemHandler) SaveUser(c *gin.Context) {
	var payload service.SaveUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.SaveUser(c.Request.Context(), operatorUserID.(uint), payload); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) DeleteUser(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.DeleteUser(c.Request.Context(), operatorUserID.(uint), userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var payload struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Status != 0 && payload.Status != 1 {
		response.Error(c, http.StatusBadRequest, "invalid status")
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.UpdateUserStatus(c.Request.Context(), operatorUserID.(uint), userID, payload.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) ListRoles(c *gin.Context) {
	var query service.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	roles, err := h.system.ListRoles(c.Request.Context(), operatorUserID.(uint), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *SystemHandler) SaveRole(c *gin.Context) {
	var payload service.SaveRolePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.SaveRole(c.Request.Context(), operatorUserID.(uint), payload); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) DeleteRole(c *gin.Context) {
	roleID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.DeleteRole(c.Request.Context(), operatorUserID.(uint), roleID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) UpdateRoleStatus(c *gin.Context) {
	roleID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var payload struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Status != 0 && payload.Status != 1 {
		response.Error(c, http.StatusBadRequest, "invalid status")
		return
	}
	operatorUserID, _ := c.Get(contextx.UserIDKey)
	if err := h.system.UpdateRoleStatus(c.Request.Context(), operatorUserID.(uint), roleID, payload.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.system.ListPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, permissions)
}

func (h *SystemHandler) SavePermission(c *gin.Context) {
	var payload service.SavePermissionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.SavePermission(c.Request.Context(), payload); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) DeletePermission(c *gin.Context) {
	permissionID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.DeletePermission(c.Request.Context(), permissionID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) ListMenus(c *gin.Context) {
	var query service.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	menus, err := h.system.ListMenus(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

func (h *SystemHandler) SaveMenu(c *gin.Context) {
	var payload model.Menu
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.SaveMenu(c.Request.Context(), payload); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) DeleteMenu(c *gin.Context) {
	menuID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.DeleteMenu(c.Request.Context(), menuID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) ListLogs(c *gin.Context) {
	var query service.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	logs, err := h.system.ListLogs(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, logs)
}

func (h *SystemHandler) GenerateOpenAPIHeaders(c *gin.Context) {
	var payload struct {
		Method string `json:"method" binding:"required"`
		Path   string `json:"path" binding:"required"`
		Body   string `json:"body"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	headers, err := utils.NewOpenAPIHeaders(h.cfg.OpenAPIKey, h.cfg.OpenAPISecret, payload.Method, payload.Path, []byte(payload.Body))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{
		"headers": gin.H{
			"X-API-Key":    headers.APIKey,
			"X-Timestamp":  headers.Timestamp,
			"X-Nonce":      headers.Nonce,
			"X-Signature":  headers.Signature,
			"Content-Type": "application/json",
		},
	})
}

func (h *SystemHandler) ListDevices(c *gin.Context) {
	var query service.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	devices, err := h.system.ListDevices(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, devices)
}

func (h *SystemHandler) UpdateDeviceStatus(c *gin.Context) {
	deviceID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var payload service.UpdateDeviceStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.UpdateDeviceStatus(c.Request.Context(), deviceID, payload.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func (h *SystemHandler) DeleteDevice(c *gin.Context) {
	deviceID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.system.DeleteDevice(c.Request.Context(), deviceID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	raw := c.Param(key)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
