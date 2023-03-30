package main

import (
	"examples/config"
	"examples/pkg/http/infra/router"
	"examples/pkg/http/infra/sql"
	"examples/pkg/http/infra/sql/engine"

	// registry の init() で使用する変数を初期化している
	// ↓ を宣言して初期化を済ませておく必要があります
	_ "examples/pkg/http/registry"
	"fmt"
	"net/http"
)

//	@title			Go Web Framework Examples.
//	@version		1.0.0
//	@description	Framework examples for Go language.

//	@host		localhost:8080
//	@BasePath	/

//	@license.name	MIT License
//	@license.url	https://github.com/ebiy0rom0/go-framework-examples/blob/develop/LICENSE

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Paste the token contained in the /signin response.

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	config.LoadConfig()
	schema := "../../db/example.db"
	if err := engine.CreateDbFile(schema); err != nil {
		return err
	}
	con, err := engine.NewMysql()
	if err != nil {
		return err
	}

	sqlh := sql.NewSqlHandler(con, con)
	router.SetRoute(sqlh)

	return http.ListenAndServe(":8080", nil)
}
