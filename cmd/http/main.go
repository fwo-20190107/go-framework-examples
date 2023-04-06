package main

import (
	"context"
	"examples/pkg/config"
	"examples/pkg/infra/cache"
	"examples/pkg/infra/framework/http/middleware"
	"examples/pkg/infra/framework/http/router"
	"examples/pkg/infra/log"
	"examples/pkg/infra/sql"
	"examples/pkg/infra/sql/engine"
	registry "examples/pkg/registry/framework/http"
	"os"
	"strconv"

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
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.InitLogger(ctx, os.Stdout)

	con, err := engine.NewMysql()
	if err != nil {
		return err
	}

	// infrastracture datasource accesssor
	sqlh := sql.NewSqlHandler(con)
	txh := sql.NewTxHandler(con)
	store := cache.NewLocalStore()

	// setup middleware
	middleware.InitAuthMiddleware(store)
	middleware.InitLoggerMiddleware(os.Stdout)

	// application DI container
	container := registry.InitializeAppContainer(sqlh, txh, store)

	router.SetRoute(container)

	return http.ListenAndServe(":"+strconv.Itoa(config.C.Server.Addr), nil)
}
