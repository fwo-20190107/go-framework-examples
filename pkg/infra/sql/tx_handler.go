package sql

import (
	"context"
	"database/sql"
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
)

type txHandler struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTxHandler(con *sql.DB) infra.TxHandler {
	return &txHandler{db: con}
}

func (h *txHandler) BeginTx(ctx context.Context, opt *sql.TxOptions) error {
	tx, err := h.db.BeginTx(ctx, opt)
	if err != nil {
		return errors.Wrap(code.CodeDatabase, err)
	}

	h.tx = tx
	return nil
}

func (h *txHandler) Commit() error {
	if err := h.tx.Commit(); err != nil {
		return errors.Wrap(code.CodeDatabase, err)
	}
	return nil
}

func (h *txHandler) Rollback() error {
	if err := h.tx.Rollback(); err != nil {
		return errors.Wrap(code.CodeDatabase, err)
	}
	return nil
}

func (h *txHandler) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := h.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}
	return result, nil
}

var _ infra.TxHandler = (*txHandler)(nil)
