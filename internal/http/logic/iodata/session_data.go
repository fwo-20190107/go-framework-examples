package iodata

type LoginInput struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}
