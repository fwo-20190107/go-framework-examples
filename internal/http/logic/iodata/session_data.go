package iodata

import validation "github.com/go-ozzo/ozzo-validation"

type SigninInput struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type SigninOutput struct {
	Token     string `json:"token"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Authority int8   `json:"authority"`
}

func (i SigninInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.LoginID, validation.Required),
		validation.Field(&i.Password, validation.Required),
	)
}
