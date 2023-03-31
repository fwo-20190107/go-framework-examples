package iodata

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SigninInput struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type SigninOutput struct {
	Token string `json:"token"`
	UserOutput
}

func (i SigninInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.LoginID, validation.Required),
		validation.Field(&i.Password, validation.Required),
	)
}
