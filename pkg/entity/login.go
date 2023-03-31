package entity

import "time"

type Login struct {
	LoginID      string
	UserID       int
	LastSignedAt *time.Time
	Password     string
}
