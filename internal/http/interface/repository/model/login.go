package model

import "time"

type Login struct {
	LoginID      string     `db:"login_id"`
	UserID       int        `db:"user_id"`
	LastSignedAt *time.Time `db:"lastSignedAt"`
	Password     string     `db:"password"`
}
