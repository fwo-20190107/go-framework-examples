package infra

import (
	"context"
	"database/sql"
)

type SqlHandler interface {
	QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
