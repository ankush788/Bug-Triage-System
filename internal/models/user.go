package models

import "time"

// User represents a registered user in the system.

type User struct {
    ID           int64     `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}