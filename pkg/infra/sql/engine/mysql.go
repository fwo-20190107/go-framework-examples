package engine

import (
	"database/sql"
	"examples/pkg/code"
	"examples/pkg/config"
	"examples/pkg/errors"

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
		ParseTime:            true,
	}
	//
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		return nil, errors.Errorf(code.CodeDatabase, err.Error())
	}
	return db, nil
}
