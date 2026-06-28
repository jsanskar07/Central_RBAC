package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	GlobalStatus string    `gorm:"default:'active'" json:"global_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Projects []ProjectUser `json:"projects,omitempty"`
	Roles    []UserRole    `json:"roles,omitempty"`
}
