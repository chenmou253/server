package service

import (
	"context"
	"fmt"

	"go-admin/server/internal/model"
	"go-admin/server/internal/repository"
	"go-admin/server/internal/utils"
)

type SystemService struct {
	repo *repository.Repository
}

const superAdminRoleCode = "super-admin"

func NewSystemService(repo *repository.Repository) *SystemService {
	return &SystemService{repo: repo}
}

type SaveRolePayload struct {
	ID            uint   `json:"id"`
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Description   string `json:"description"`
	Status        int8   `json:"status"`
	PermissionIDs []uint `json:"permissionIds"`
	MenuIDs       []uint `json:"menuIds"`
}

type SaveUserPayload struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

type SavePermissionPayload struct {
	ID     uint   `json:"id"`
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Module string `json:"module" binding:"required"`
}

type UpdateDeviceStatusPayload struct {
	Status uint8 `json:"status"`
}

type PageQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type PageResult struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func (q PageQuery) Normalize() PageQuery {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	if q.PageSize > 500 {
		q.PageSize = 500
	}
	return q
}

func (s *SystemService) ListUsers(ctx context.Context, operatorUserID uint, query PageQuery) (*PageResult, error) {
	query = query.Normalize()
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return nil, err
	}
	users, total, err := s.repo.ListUsers(ctx, query.Page, query.PageSize, !isSuperAdmin)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: users, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SystemService) SaveUser(ctx context.Context, operatorUserID uint, payload SaveUserPayload) error {
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return err
	}
	var user *model.User
	if payload.ID > 0 {
		current, err := s.repo.GetUserByID(ctx, payload.ID)
		if err != nil {
			return err
		}
		if !isSuperAdmin && hasRoleCode(current.Roles, superAdminRoleCode) {
			return fmt.Errorf("cannot operate super admin user")
		}
		user = current
	} else {
		user = &model.User{}
	}

	user.Username = payload.Username
	user.Nickname = payload.Nickname
	user.Status = payload.Status
	if user.Status == 0 {
		user.Status = 1
	}

	if payload.Password != "" {
		password, err := utils.HashPassword(payload.Password)
		if err != nil {
			return err
		}
		user.Password = password
	}
	if payload.ID > 0 {
		if err := s.repo.SaveUser(ctx, user); err != nil {
			return err
		}
	} else {
		if user.Password == "" {
			password, err := utils.HashPassword("123456")
			if err != nil {
				return err
			}
			user.Password = password
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return err
		}
	}
	roles, err := s.repo.GetRolesByIDs(ctx, payload.RoleIDs)
	if err != nil {
		return err
	}
	if !isSuperAdmin && hasRoleCode(roles, superAdminRoleCode) {
		return fmt.Errorf("cannot assign super admin role")
	}
	return s.repo.ReplaceUserRoles(ctx, user, roles)
}

func (s *SystemService) DeleteUser(ctx context.Context, operatorUserID uint, userID uint) error {
	if operatorUserID == userID {
		return fmt.Errorf("cannot delete current user")
	}
	return s.repo.DeleteUser(ctx, userID)
}

func (s *SystemService) UpdateUserStatus(ctx context.Context, operatorUserID uint, userID uint, status int8) error {
	if operatorUserID == userID {
		return fmt.Errorf("cannot change current user status")
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Status = status
	return s.repo.SaveUser(ctx, user)
}

func (s *SystemService) ListRoles(ctx context.Context, operatorUserID uint, query PageQuery) (*PageResult, error) {
	query = query.Normalize()
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return nil, err
	}
	roles, total, err := s.repo.ListRoles(ctx, query.Page, query.PageSize, !isSuperAdmin)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: roles, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SystemService) SaveRole(ctx context.Context, operatorUserID uint, payload SaveRolePayload) error {
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return err
	}
	if !isSuperAdmin && payload.Code == superAdminRoleCode {
		return fmt.Errorf("cannot operate super admin role")
	}
	if payload.ID > 0 {
		current, err := s.repo.GetRoleByID(ctx, payload.ID)
		if err != nil {
			return err
		}
		if !isSuperAdmin && current.Code == superAdminRoleCode {
			return fmt.Errorf("cannot operate super admin role")
		}
	}
	role := &model.Role{
		BaseModel:   model.BaseModel{ID: payload.ID},
		Name:        payload.Name,
		Code:        payload.Code,
		Description: payload.Description,
		Status:      payload.Status,
	}
	if role.Status == 0 {
		role.Status = 1
	}
	if payload.ID > 0 {
		if err := s.repo.SaveRole(ctx, role); err != nil {
			return err
		}
	} else {
		if err := s.repo.CreateRole(ctx, role); err != nil {
			return err
		}
	}
	permissions, err := s.repo.GetPermissionsByIDs(ctx, payload.PermissionIDs)
	if err != nil {
		return err
	}
	menus, err := s.repo.GetMenusByIDs(ctx, payload.MenuIDs)
	if err != nil {
		return err
	}
	if err := s.repo.ReplaceRolePermissions(ctx, role, permissions); err != nil {
		return err
	}
	return s.repo.ReplaceRoleMenus(ctx, role, menus)
}

func (s *SystemService) DeleteRole(ctx context.Context, operatorUserID uint, roleID uint) error {
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return err
	}
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if !isSuperAdmin && role.Code == superAdminRoleCode {
		return fmt.Errorf("cannot operate super admin role")
	}
	usedCount, err := s.repo.CountUsersByRoleID(ctx, roleID)
	if err != nil {
		return err
	}
	if usedCount > 0 {
		return fmt.Errorf("role is assigned to users and cannot be deleted")
	}
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *SystemService) UpdateRoleStatus(ctx context.Context, operatorUserID uint, roleID uint, status int8) error {
	isSuperAdmin, err := s.isSuperAdmin(ctx, operatorUserID)
	if err != nil {
		return err
	}
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if !isSuperAdmin && role.Code == superAdminRoleCode {
		return fmt.Errorf("cannot operate super admin role")
	}
	role.Status = status
	return s.repo.SaveRole(ctx, role)
}

func (s *SystemService) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *SystemService) SavePermission(ctx context.Context, payload SavePermissionPayload) error {
	permission := &model.Permission{
		BaseModel: model.BaseModel{ID: payload.ID},
		Name:      payload.Name,
		Code:      payload.Code,
		Module:    payload.Module,
	}
	if payload.ID > 0 {
		return s.repo.SavePermission(ctx, permission)
	}
	return s.repo.CreatePermission(ctx, permission)
}

func (s *SystemService) DeletePermission(ctx context.Context, permissionID uint) error {
	return s.repo.DeletePermission(ctx, permissionID)
}

func (s *SystemService) ListMenus(ctx context.Context, query PageQuery) (*PageResult, error) {
	query = query.Normalize()
	menus, total, err := s.repo.ListMenus(ctx, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: menus, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SystemService) SaveMenu(ctx context.Context, payload model.Menu) error {
	if payload.ID > 0 {
		return s.repo.SaveMenu(ctx, &payload)
	}
	return s.repo.CreateMenu(ctx, &payload)
}

func (s *SystemService) DeleteMenu(ctx context.Context, menuID uint) error {
	childCount, err := s.repo.CountMenusByParentID(ctx, menuID)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("menu has child items and cannot be deleted")
	}
	return s.repo.DeleteMenu(ctx, menuID)
}

func (s *SystemService) ListLogs(ctx context.Context, query PageQuery) (*PageResult, error) {
	query = query.Normalize()
	logs, total, err := s.repo.ListLogs(ctx, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: logs, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SystemService) ListDevices(ctx context.Context, query PageQuery) (*PageResult, error) {
	query = query.Normalize()
	devices, total, err := s.repo.ListDevices(ctx, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}
	return &PageResult{List: devices, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SystemService) UpdateDeviceStatus(ctx context.Context, deviceID uint, status uint8) error {
	if status != 0 && status != 1 {
		return fmt.Errorf("invalid status")
	}
	device, err := s.repo.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return err
	}
	device.Status = status
	return s.repo.SaveDevice(ctx, device)
}

func (s *SystemService) DeleteDevice(ctx context.Context, deviceID uint) error {
	return s.repo.DeleteDevice(ctx, deviceID)
}

func (s *SystemService) isSuperAdmin(ctx context.Context, userID uint) (bool, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return hasRoleCode(user.Roles, superAdminRoleCode), nil
}

func hasRoleCode(roles []model.Role, code string) bool {
	for _, role := range roles {
		if role.Code == code {
			return true
		}
	}
	return false
}
