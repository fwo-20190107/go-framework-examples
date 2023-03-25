package iodata

import (
	customValidaton "examples/internal/http/logic/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SignupInput struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (i SignupInput) Validate() error {
	err := validation.ValidateStruct(&i,
		validation.Field(&i.LoginID, validation.Required, validation.Length(4, 16), customValidaton.NgWord),
		validation.Field(&i.Password, validation.Required, validation.Length(8, 32)),
	)
	if err != nil {
		return err
	}
	return nil
}
