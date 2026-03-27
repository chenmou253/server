package repository

import (
	"context"
	"strings"

	"go-admin/server/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

const superAdminRoleCode = "super-admin"

func New(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).
		Preload("Roles.Permissions").
		Preload("Roles.Menus").
		Where("username = ?", username).
		First(&user).Error
	return &user, err
}

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).
		Preload("Roles.Permissions").
		Preload("Roles.Menus").
		First(&user, userID).Error
	return &user, err
}

func (r *Repository) ListUsers(ctx context.Context, page int, pageSize int, excludeSuperAdmin bool) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.User{})
	if excludeSuperAdmin {
		db = db.Where(`NOT EXISTS (
			SELECT 1 FROM user_roles
			JOIN roles ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = users.id AND roles.code = ?
		)`, superAdminRoleCode)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := r.DB.WithContext(ctx).
		Preload("Roles").
		Order("id desc")
	if excludeSuperAdmin {
		query = query.Where(`NOT EXISTS (
			SELECT 1 FROM user_roles
			JOIN roles ON roles.id = user_roles.role_id
			WHERE user_roles.user_id = users.id AND roles.code = ?
		)`, superAdminRoleCode)
	}
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *Repository) SaveUser(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint) error {
	return r.DB.WithContext(ctx).Delete(&model.User{}, userID).Error
}

func (r *Repository) ReplaceUserRoles(ctx context.Context, user *model.User, roles []model.Role) error {
	return r.DB.WithContext(ctx).Model(user).Association("Roles").Replace(roles)
}

func (r *Repository) ListRoles(ctx context.Context, page int, pageSize int, excludeSuperAdmin bool) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.Role{})
	if excludeSuperAdmin {
		db = db.Where("code <> ?", superAdminRoleCode)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := r.DB.WithContext(ctx).
		Preload("Permissions").
		Preload("Menus").
		Order("id asc")
	if excludeSuperAdmin {
		query = query.Where("code <> ?", superAdminRoleCode)
	}
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

func (r *Repository) GetRolesByIDs(ctx context.Context, ids []uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.WithContext(ctx).Find(&roles, ids).Error
	return roles, err
}

func (r *Repository) SaveRole(ctx context.Context, role *model.Role) error {
	return r.DB.WithContext(ctx).Save(role).Error
}

func (r *Repository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.DB.WithContext(ctx).Create(role).Error
}

func (r *Repository) GetRoleByID(ctx context.Context, roleID uint) (*model.Role, error) {
	var role model.Role
	err := r.DB.WithContext(ctx).First(&role, roleID).Error
	return &role, err
}

func (r *Repository) DeleteRole(ctx context.Context, roleID uint) error {
	return r.DB.WithContext(ctx).Delete(&model.Role{}, roleID).Error
}

func (r *Repository) CountUsersByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Table("user_roles").Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

func (r *Repository) ReplaceRolePermissions(ctx context.Context, role *model.Role, permissions []model.Permission) error {
	return r.DB.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

func (r *Repository) ReplaceRoleMenus(ctx context.Context, role *model.Role, menus []model.Menu) error {
	return r.DB.WithContext(ctx).Model(role).Association("Menus").Replace(menus)
}

func (r *Repository) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.WithContext(ctx).Order("module asc, id asc").Find(&permissions).Error
	return permissions, err
}

func (r *Repository) GetPermissionsByIDs(ctx context.Context, ids []uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.WithContext(ctx).Find(&permissions, ids).Error
	return permissions, err
}

func (r *Repository) SavePermission(ctx context.Context, permission *model.Permission) error {
	return r.DB.WithContext(ctx).Save(permission).Error
}

func (r *Repository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return r.DB.WithContext(ctx).Create(permission).Error
}

func (r *Repository) DeletePermission(ctx context.Context, permissionID uint) error {
	return r.DB.WithContext(ctx).Delete(&model.Permission{}, permissionID).Error
}

func (r *Repository) ListMenus(ctx context.Context, page int, pageSize int) ([]model.Menu, int64, error) {
	var menus []model.Menu
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.Menu{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.DB.WithContext(ctx).
		Order("sort asc, id asc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&menus).Error
	return menus, total, err
}

func (r *Repository) GetMenusByIDs(ctx context.Context, ids []uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.DB.WithContext(ctx).Find(&menus, ids).Error
	return menus, err
}

func (r *Repository) SaveMenu(ctx context.Context, menu *model.Menu) error {
	return r.DB.WithContext(ctx).Save(menu).Error
}

func (r *Repository) CreateMenu(ctx context.Context, menu *model.Menu) error {
	return r.DB.WithContext(ctx).Create(menu).Error
}

func (r *Repository) DeleteMenu(ctx context.Context, menuID uint) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Table("role_menus").Where("menu_id = ?", menuID).Delete(nil).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Delete(&model.Menu{}, menuID).Error
	})
}

func (r *Repository) CountMenusByParentID(ctx context.Context, parentID uint) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.Menu{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

func (r *Repository) CreateLog(ctx context.Context, log *model.OperationLog) error {
	return r.DB.WithContext(ctx).Create(log).Error
}

func (r *Repository) ListLogs(ctx context.Context, page int, pageSize int) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.OperationLog{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.DB.WithContext(ctx).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error
	return logs, total, err
}

func (r *Repository) ListDevices(ctx context.Context, page int, pageSize int) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.Device{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.DB.WithContext(ctx).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&devices).Error
	return devices, total, err
}

func (r *Repository) GetDeviceByID(ctx context.Context, deviceID uint) (*model.Device, error) {
	var device model.Device
	err := r.DB.WithContext(ctx).First(&device, deviceID).Error
	return &device, err
}

func (r *Repository) SaveDevice(ctx context.Context, device *model.Device) error {
	return r.DB.WithContext(ctx).Save(device).Error
}

func (r *Repository) DeleteDevice(ctx context.Context, deviceID uint) error {
	return r.DB.WithContext(ctx).Delete(&model.Device{}, deviceID).Error
}

func (r *Repository) CreateDevice(ctx context.Context, device *model.Device) error {
	return r.DB.WithContext(ctx).Create(device).Error
}

type DeviceQuery struct {
	DeviceNo   string
	MerchantID string
	Status     *uint8
}

func (r *Repository) ListDevicesByQuery(ctx context.Context, page int, pageSize int, query DeviceQuery) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64
	db := r.DB.WithContext(ctx).Model(&model.Device{})
	db = applyDeviceQuery(db, query)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := applyDeviceQuery(r.DB.WithContext(ctx).Model(&model.Device{}), query).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&devices).Error
	return devices, total, err
}

func (r *Repository) GetDeviceByDeviceNo(ctx context.Context, deviceNo string) (*model.Device, error) {
	var device model.Device
	err := r.DB.WithContext(ctx).Where("device_no = ?", deviceNo).First(&device).Error
	return &device, err
}

func applyDeviceQuery(db *gorm.DB, query DeviceQuery) *gorm.DB {
	if query.DeviceNo != "" {
		db = db.Where("device_no LIKE ?", "%"+strings.TrimSpace(query.DeviceNo)+"%")
	}
	if query.MerchantID != "" {
		db = db.Where("merchant_id = ?", strings.TrimSpace(query.MerchantID))
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	return db
}
