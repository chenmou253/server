# Go Admin Server

Gin + GORM + MySQL + Redis 的后台服务。

## 功能

- JWT 登录、刷新 Token、退出登录
- RBAC：角色 -> 权限字符串，不使用 Casbin
- 菜单管理、角色管理、用户管理
- 操作日志中间件统一记录
- Redis 存储刷新 Token

## 启动

1. 复制 `.env.example` 为 `.env`
2. 创建 MySQL 数据库 `go_admin`
3. 执行 `go mod tidy`
4. 执行 `go run ./cmd/server`

默认账号：`admin / 123456`
