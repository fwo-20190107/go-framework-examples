package sql

import (
	"context"
	"database/sql"
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"

	"github.com/tanimutomo/sqlfile"
)

type sqlHandler struct {
	con *sql.DB
}

func NewSqlHandler(con *sql.DB) *sqlHandler {
	return &sqlHandler{
		con: con,
	}
}

func (h *sqlHandler) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	stmt, err := h.con.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	return row, nil
}

func (h *sqlHandler) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	stmt, err := h.con.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(code.CodeDatabase, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		c := code.CodeDatabase
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c = code.CodeNotFound
		}
		return nil, errors.Wrap(c, err)
	}
	return rows, nil
}

func InitializeDb(con *sql.DB) error {
	s := sqlfile.New()
	if err := s.Directory("../testdata"); err != nil {
		return errors.Wrap(code.CodeInternal, err)
	}
	if _, err := s.Exec(con); err != nil {
		return errors.Wrap(code.CodeDatabase, err)
	}
	return nil
}

var _ infra.SqlHandler = (*sqlHandler)(nil)
