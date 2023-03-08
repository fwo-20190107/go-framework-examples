package message

import "errors"

var (
	ErrParseForm       = errors.New("failed parse form")
	ErrLoginNotFound   = errors.New("loginID is not exist")
	ErrUnmatchPassword = errors.New("password does not match")
)
