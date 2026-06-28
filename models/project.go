package models

import (
	"time"
)

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	APIKey      string    `gorm:"uniqueIndex" json:"api_key,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectUser struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ProjectID uint      `gorm:"primaryKey" json:"project_id"`
	Status    string    `gorm:"default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
}
