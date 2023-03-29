package entity

import "time"

type Login struct {
	LoginID      string     `json:"login_id"`
	UserID       int        `json:"user_id"`
	LastSignedAt *time.Time `json:"last_signed_at"`
	Password     string     `json:"password"`
}
