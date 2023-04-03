package infra

import (
	"context"
	"database/sql"
)

type Executer interface {
	Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TxHandler interface {
	Executer
	BeginTx(ctx context.Context, opt *sql.TxOptions) error
	Commit() error
	Rollback() error
}
