package model

type User struct {
	UserID    int    `db:"user_id"`
	Name      string `db:"name"`
	Authority int8   `db:"authority"`
}

type Login struct {
	LoginID  string `db:"login_id"`
	UserID   int    `db:"user_id"`
	Password string `db:"password"`
}
