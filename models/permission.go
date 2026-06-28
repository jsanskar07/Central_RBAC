package models

import (
	"time"
)

type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index" json:"project_id"` // Can be null if it's a global/standard permission
	Resource  string    `gorm:"not null" json:"resource"`
	Action    string    `gorm:"not null" json:"action"`
	CreatedAt time.Time `json:"created_at"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"-"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`

	Role       Role       `gorm:"foreignKey:RoleID" json:"-"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"-"`
}
