package infra

import (
	"context"
	"database/sql"
)

type Executor interface {
	Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TxHandler interface {
	Executor
	BeginTx(ctx context.Context, opt *sql.TxOptions) error
	Commit() error
	Rollback() error
}
