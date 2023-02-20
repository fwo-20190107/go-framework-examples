package model

type User struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Authority int8   `json:"authority"`
}

type Login struct {
	LoginID  string `json:"login_id"`
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}
