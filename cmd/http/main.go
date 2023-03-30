package main

import (
	"examples/pkg/config"
	"examples/pkg/infra/router"
	"examples/pkg/infra/sql"
	"examples/pkg/infra/sql/engine"

	// registry の init() で使用する変数を初期化している
	// ↓ を宣言して初期化を済ませておく必要があります
	_ "examples/pkg/logger"
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

	con, err := engine.NewMysql()
	if err != nil {
		return err
	}

	sqlh := sql.NewSqlHandler(con, con)
	router.SetRoute(sqlh)

	return http.ListenAndServe(":8080", nil)
}
