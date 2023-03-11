package main

import (
	"examples/internal/http/infra/router"
	"examples/internal/http/infra/sql"
	"examples/internal/http/infra/sql/engine"

	// registry の init() で使用する変数を初期化している
	// ↓ を宣言して初期化を済ませておく必要があります
	_ "examples/internal/http/registry"
	"fmt"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	schema := "../../db/example.db"
	if err := engine.CreateDbFile(schema); err != nil {
		return err
	}
	con, err := engine.NewSqlite3(schema)
	if err != nil {
		return err
	}

	if err := sql.InitializeDb(con); err != nil {
		return err
	}

	sqlh := sql.NewSqlHandler(con, con)
	router.SetRoute(sqlh)

	return http.ListenAndServe(":8080", nil)
}
