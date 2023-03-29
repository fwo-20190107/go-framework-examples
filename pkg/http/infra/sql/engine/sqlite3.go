package engine

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite3(schema string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", schema)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateDbFile(schema string) error {
	dir := filepath.Dir(schema)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0766); err != nil {
			return err
		}
	}
	if _, err := os.Create(schema); err != nil {
		return err
	}
	return nil
}
