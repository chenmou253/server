package model

import "time"

type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Nickname string `gorm:"size:64;not null" json:"nickname"`
	Password string `gorm:"size:255;not null" json:"-"`
	Status   int8   `gorm:"default:1" json:"status"`
	Roles    []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

type Role struct {
	BaseModel
	Name        string       `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string       `gorm:"size:255" json:"description"`
	Status      int8         `gorm:"default:1" json:"status"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
}

type Permission struct {
	BaseModel
	Name   string `gorm:"size:64;not null" json:"name"`
	Code   string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Module string `gorm:"size:64;not null" json:"module"`
}

type Menu struct {
	BaseModel
	ParentID   uint   `gorm:"default:0;index" json:"parentId"`
	Name       string `gorm:"size:64;not null" json:"name"`
	Title      string `gorm:"size:64;not null" json:"title"`
	Path       string `gorm:"size:128;not null" json:"path"`
	Component  string `gorm:"size:128;not null" json:"component"`
	Icon       string `gorm:"size:64" json:"icon"`
	MenuType   string `gorm:"size:16;not null;default:menu" json:"menuType"`
	Permission string `gorm:"size:64" json:"permission"`
	Sort       int    `gorm:"default:0" json:"sort"`
	Hidden     bool   `gorm:"default:false" json:"hidden"`
	KeepAlive  bool   `gorm:"default:false" json:"keepAlive"`
	Status     int8   `gorm:"default:1" json:"status"`
}

type Device struct {
	ID         uint   `gorm:"primaryKey;column:id" json:"id"`
	DeviceNo   string `gorm:"column:device_no" json:"deviceNo"`
	MerchantID string `gorm:"column:merchant_id" json:"merchantId"`
	Status     uint8  `gorm:"column:status" json:"status" comment:"设备状态0待机1接单"`
	IP         string `gorm:"column:ip" json:"ip"`
	CreateT    uint64 `gorm:"column:create_t" json:"createT"`
}

func (Device) TableName() string {
	return "device"
}

type OperationLog struct {
	BaseModel
	UserID     uint   `gorm:"index" json:"userId"`
	Username   string `gorm:"size:64" json:"username"`
	Method     string `gorm:"size:16" json:"method"`
	Path       string `gorm:"size:255" json:"path"`
	StatusCode int    `json:"statusCode"`
	IP         string `gorm:"size:64" json:"ip"`
	UserAgent  string `gorm:"size:255" json:"userAgent"`
	Action     string `gorm:"size:255" json:"action"`
	LatencyMS  int64  `json:"latencyMs"`
	RequestID  string `gorm:"size:64" json:"requestId"`
}
