package infra

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tanimutomo/sqlfile"
)

func InitializeDb(schema string) error {
	dir := filepath.Dir(schema)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	if _, err := os.Create(schema); err != nil {
		return err
	}

	con, err := NewConnection(schema)
	if err != nil {
		return err
	}

	s := sqlfile.New()
	if err := s.Directory("testdata"); err != nil {
		return err
	}
	if _, err := s.Exec(con); err != nil {
		return err
	}
	return nil
}

func NewConnection(schema string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", schema)
	if err != nil {
		return nil, err
	}
	return db, nil
}
