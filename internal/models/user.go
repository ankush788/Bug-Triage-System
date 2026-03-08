package models

import (
	"time"
)

// User represents a registered user in the system.
type User struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email        string         `gorm:"column:email;unique;not null" json:"email"`
	PasswordHash string         `gorm:"column:password_hash;not null" json:"-"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`

}