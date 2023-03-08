package message

import "errors"

var (
	ErrParseForm       = errors.New("failed parse form")
	ErrUserNotFound    = errors.New("loginID is not exist")
	ErrUnmatchPassword = errors.New("password does not match")
)
