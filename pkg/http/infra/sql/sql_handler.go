package sql

import (
	"context"
	"database/sql"
	"examples/code"
	"examples/errors"
	"examples/pkg/http/adapter/infra"

	"github.com/tanimutomo/sqlfile"
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
		return nil, errors.Wrap(code.ErrDatabase, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(code.ErrDatabase, err)
	}
	return res, nil
}

func (h *sqlHandler) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	stmt, err := h.slave.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(code.ErrDatabase, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	return row, nil
}

func (h *sqlHandler) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	stmt, err := h.slave.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(code.ErrDatabase, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		c := code.ErrDatabase
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c = code.ErrNotFound
		}
		return nil, errors.Wrap(c, err)
	}
	return rows, nil
}

func InitializeDb(con *sql.DB) error {
	s := sqlfile.New()
	if err := s.Directory("../testdata"); err != nil {
		return errors.Wrap(code.ErrInternal, err)
	}
	if _, err := s.Exec(con); err != nil {
		return errors.Wrap(code.ErrDatabase, err)
	}
	return nil
}

var _ infra.SqlHandler = (*sqlHandler)(nil)
