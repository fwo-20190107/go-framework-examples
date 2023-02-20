package infra

import (
	"context"
	"database/sql"
	"examples/interface/repository"
)

type sqlHandler struct {
	master *sql.DB
	slave  *sql.DB
}

func NewSqlHandler(master, slave *sql.DB) *sqlHandler {
	return &sqlHandler{
		master: master,
		slave:  slave,
	}
}

func (h *sqlHandler) Execute(ctx context.Context, query string, args ...any) (sql.Result, error) {
	stmt, err := h.master.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *sqlHandler) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	stmt, err := h.slave.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	return row, nil
}

func (h *sqlHandler) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	stmt, err := h.slave.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

var _ repository.SqlHandler = (*sqlHandler)(nil)
