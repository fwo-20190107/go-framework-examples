package main

import (
	"examples/pkg/config"
	"examples/pkg/infra/cache"
	"examples/pkg/infra/framework/gin/middleware"
	"examples/pkg/infra/framework/gin/router"
	"examples/pkg/infra/sql"
	"examples/pkg/infra/sql/engine"
	registry "examples/pkg/registry/framework/gin"
	"fmt"
	"os"
)

// https://github.com/gin-gonic/gin/issues/1734

// pathパラメータの型指定を許容しないのは、単純に複雑すぎるからとのこと
// (context).Param("key")はstringオンリーのため適宜型変換を行う事
// https://github.com/gin-gonic/gin/issues/846#issuecomment-439312131
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

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

	// application DI container
	container := registry.InitializeAppContainer(sqlh, txh, store)

	r := router.SetRoute(container)

	return r.Run()
}
