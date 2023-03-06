package main

import (
	"examples/internal/http/infra/routes"
	"examples/internal/http/infra/sql"
	"examples/internal/http/infra/sql/driver"
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
	con, err := driver.NewConnection("db/example.db")
	if err != nil {
		return err
	}

	sqlh := sql.NewSqlHandler(con, con)
	routes.SetRoute(sqlh)

	return http.ListenAndServe(":8080", nil)
}
