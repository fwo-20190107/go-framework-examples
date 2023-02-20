package framework

import (
	"examples/infra"
	"fmt"
	"net/http"
)

type Framework string

var (
	FwGin     Framework = "gin"
	FwEcho    Framework = "echo"
	FwGoa     Framework = "goa"
	FwGorilla Framework = "gorilla"
)

func Start(fw Framework) error {
	fmt.Printf("select %s framework.\n", fw)

	con, err := infra.NewConnection("db/example.db")
	if err != nil {
		return err
	}

	//
	sqlh := infra.NewSqlHandler(con, con)

	switch fw {
	case FwGin:
		// todo: new gin
	case FwEcho:
		// todo: new echo
	case FwGoa:
		// todo: new goa
	case FwGorilla:
		// todo: new gorilla/mux
	default:
		New(sqlh)
		return http.ListenAndServe(":8080", nil)
	}
	return nil
}
