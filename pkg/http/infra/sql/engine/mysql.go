package engine

import (
	"database/sql"
	"examples/code"
	"examples/config"
	"examples/errors"

	"github.com/go-sql-driver/mysql"
)

func NewMysql() (*sql.DB, error) {
	cnf := &mysql.Config{
		User:                 config.C.DB.User,
		Passwd:               config.C.DB.Passwd,
		Net:                  config.C.DB.Net,
		Addr:                 config.C.DB.Addr,
		DBName:               config.C.DB.DBName,
		AllowNativePasswords: true,
	}
	//
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		return nil, errors.Errorf(code.ErrDatabase, err.Error())
	}
	return db, nil
}
