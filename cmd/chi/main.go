package main

import (
	"context"
	"examples/pkg/config"
	"examples/pkg/infra/cache"
	"examples/pkg/infra/framework/chi/middleware"
	"examples/pkg/infra/framework/chi/router"
	"examples/pkg/infra/log"
	"examples/pkg/infra/sql"
	"examples/pkg/infra/sql/engine"
	registry "examples/pkg/registry/framework/chi"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

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
	// middleware.InitLoggerMiddleware(os.Stdout)

	// application DI container
	container := registry.InitializeAppContainer(sqlh, txh, store)

	r := router.SetRoute(container)

	return http.ListenAndServe(":"+strconv.Itoa(config.C.Server.Addr), r)
}
