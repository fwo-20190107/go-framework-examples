package model

type User struct {
	UserID    int    `db:"user_id"`
	Name      string `db:"name"`
	Authority int8   `db:"authority"`
}
