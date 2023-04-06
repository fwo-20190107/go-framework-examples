package main

import (
	"context"
	"examples/pkg/config"
	"examples/pkg/infra/cache"
	"examples/pkg/infra/framework/gin/middleware"
	"examples/pkg/infra/framework/gin/router"
	"examples/pkg/infra/log"
	"examples/pkg/infra/sql"
	"examples/pkg/infra/sql/engine"
	registry "examples/pkg/registry/framework/gin"
	"fmt"
	"os"
	"strconv"
)

// 通常の context.Background() で生成される context との共存は難しいので、gin.Context を使えば良い
// 実装を見ると分かるが、gin.Context は context.Context を満たす構造体である(実質的なスーパーセット)
// https://github.com/gin-gonic/gin/issues/1734

// ginのルーティングはメソッドの制限とpathパラメータの名前付けに限定される
// 以下のような問題もあり、これが致命的な場合は選択肢から外す必要がある
// https://www.irohabook.com/gin-router

// pathパラメータ取得は(context).Param("key")を使う
// (context).Param("key")はstringオンリーのため適宜型変換を行う必要がある
// (context).BindUri(&struct)も使えるようになってるらしい
// https://github.com/gin-gonic/gin/issues/846

// pathパラメータの型指定を許容しないのは、単純に複雑すぎるから
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

	// application DI container
	container := registry.InitializeAppContainer(sqlh, txh, store)

	r := router.SetRoute(container)

	return r.Run(":" + strconv.Itoa(config.C.Server.Addr))
}
