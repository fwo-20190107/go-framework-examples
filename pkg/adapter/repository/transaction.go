package repository

import (
	"context"
	"database/sql"
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logic/repository"
)

type ctxkey struct{}

type tx struct {
	infra.TxHandler
}

var txKey = ctxkey{}

func NewTransaction(handler infra.TxHandler) repository.Transaction {
	return &tx{handler}
}

func (t *tx) Do(ctx context.Context, fn repository.TxExecFunc) (interface{}, error) {
	opt := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	if err := t.BeginTx(ctx, opt); err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}
	defer t.Rollback()

	ctx = context.WithValue(ctx, &txKey, t)
	v, err := fn(ctx)
	if err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}

	if err := t.Commit(); err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}
	return v, nil
}

func getExecutor(ctx context.Context) (infra.Executor, error) {
	e, ok := ctx.Value(&txKey).(infra.Executor)
	if !ok {
		return nil, errors.Errorf(code.CodeImplements, "failed to load executor from context")
	}
	return e, nil
}

var _ repository.Transaction = (*tx)(nil)
