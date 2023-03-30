package code

type ErrorCode string

var (
	ErrOK ErrorCode = "ok"

	ErrBadRequest   ErrorCode = "bad request"
	ErrUnauthorized ErrorCode = "unauthorized"
	ErrValidParam   ErrorCode = "valid parameter"
	ErrNotFound     ErrorCode = "not found"
	ErrOutOfTerm    ErrorCode = "out of term"

	ErrDatabase ErrorCode = "database error"
	ErrInternal ErrorCode = "internal error"

	ErrUnknown ErrorCode = "unknown error"
)

func (e ErrorCode) Error() string {
	return string(e)
}

var _ error = (*ErrorCode)(nil)
