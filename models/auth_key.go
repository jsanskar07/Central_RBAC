package models

import (
	"time"
)

type AuthKey struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PrivateKey string    `gorm:"not null" json:"-"`
	PublicKey  string    `gorm:"not null" json:"public_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
