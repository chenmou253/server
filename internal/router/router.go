package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-admin/server/internal/config"
	"go-admin/server/internal/handler"
	"go-admin/server/internal/middleware"
	"go-admin/server/internal/repository"
	"go-admin/server/internal/service"
)

func New(cfg *config.Config, repo *repository.Repository, authHandler *handler.AuthHandler, systemHandler *handler.SystemHandler, publicHandler *handler.PublicHandler, rbac *service.RBACService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(repo))
	r.StaticFile("/api-docs/swagger.yaml", "docs/swag/swagger.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api-docs/swagger.yaml")))

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)

		open := api.Group("/open")
		open.Use(middleware.OpenAPIAuth(cfg))
		open.GET("/devices", publicHandler.ListDevices)
		open.POST("/devices", publicHandler.RegisterDevice)
		open.GET("/devices/:deviceNo", publicHandler.GetDevice)
		open.POST("/devices/:deviceNo/status", publicHandler.UpdateDeviceStatus)

		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg, rbac))
		protected.GET("/auth/profile", authHandler.Profile)
		protected.POST("/auth/logout", authHandler.Logout)

		system := protected.Group("/system")
		system.GET("/menus", systemHandler.ListMenus)
		system.POST("/menus", middleware.RequirePermission("menu:save"), systemHandler.SaveMenu)
		system.DELETE("/menus/:id", middleware.RequirePermission("menu:save"), systemHandler.DeleteMenu)
		system.GET("/users", middleware.RequirePermission("user:list"), systemHandler.ListUsers)
		system.POST("/users", middleware.RequirePermission("user:save"), systemHandler.SaveUser)
		system.DELETE("/users/:id", middleware.RequirePermission("user:delete"), systemHandler.DeleteUser)
		system.PATCH("/users/:id/status", middleware.RequirePermission("user:status"), systemHandler.UpdateUserStatus)
		system.GET("/roles", middleware.RequirePermission("role:list"), systemHandler.ListRoles)
		system.POST("/roles", middleware.RequirePermission("role:save"), systemHandler.SaveRole)
		system.DELETE("/roles/:id", middleware.RequirePermission("role:delete"), systemHandler.DeleteRole)
		system.PATCH("/roles/:id/status", middleware.RequirePermission("role:status"), systemHandler.UpdateRoleStatus)
		system.GET("/permissions", systemHandler.ListPermissions)
		system.POST("/permissions", middleware.RequirePermission("menu:save"), systemHandler.SavePermission)
		system.DELETE("/permissions/:id", middleware.RequirePermission("menu:save"), systemHandler.DeletePermission)
		system.GET("/logs", middleware.RequirePermission("log:list"), systemHandler.ListLogs)
		system.POST("/openapi/sign", middleware.RequirePermission("menu:save"), systemHandler.GenerateOpenAPIHeaders)
		system.GET("/devices", middleware.RequirePermission("device:list"), systemHandler.ListDevices)
		system.PATCH("/devices/:id/status", middleware.RequirePermission("device:status"), systemHandler.UpdateDeviceStatus)
		system.DELETE("/devices/:id", middleware.RequirePermission("device:delete"), systemHandler.DeleteDevice)
	}

	return r
}
