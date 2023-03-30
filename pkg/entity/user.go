package entity

type User struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Authority int8   `json:"authority"`
}
