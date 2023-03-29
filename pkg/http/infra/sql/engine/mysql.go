package engine

import (
	"database/sql"
	"examples/code"
	"examples/errors"

	"github.com/go-sql-driver/mysql"
)

func NewMysql() (*sql.DB, error) {
	cnf := &mysql.Config{
		User:                 "user",
		Passwd:               "userpass",
		Net:                  "tcp",
		Addr:                 "mysqldb010:3306",
		DBName:               "app",
		AllowNativePasswords: true,
	}
	//
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		return nil, errors.Errorf(code.ErrDatabase, err.Error())
	}
	return db, nil
}
