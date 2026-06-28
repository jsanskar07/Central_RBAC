package models

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Project     Project          `gorm:"foreignKey:ProjectID" json:"-"`
	Permissions []RolePermission `json:"permissions,omitempty"`
}

// UserRole links a user to a specific role within a project context
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Role Role `gorm:"foreignKey:RoleID" json:"-"`
}
