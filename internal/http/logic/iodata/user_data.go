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
	return validation.ValidateStruct(&i,
		validation.Field(&i.LoginID, validation.Required, validation.Length(4, 16), customValidaton.NgWord),
		validation.Field(&i.Password, validation.Required, validation.Length(8, 32)),
	)
}

type ModifyNameInput struct {
	Name string `json:"name"`
}

type ModifyAuthorityInput struct {
	Authority int8 `json:"authority"`
}

type UserModifyOutput struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Authority int8   `json:"authority"`
}

func (i ModifyNameInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Length(4, 16), customValidaton.NgWord),
	)
}

func (i ModifyAuthorityInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Authority, validation.Required, validation.Max(99)),
	)
}
