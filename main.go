package main

import (
	"examples/framework"
	"examples/infra"
	_ "examples/infra/middleware"
	"flag"
	"fmt"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	var fw framework.Framework
	flag.StringVar((*string)(&fw), "framework", "", "exec framework")
	flag.StringVar((*string)(&fw), "f", "http", "exec framework")

	if err := infra.InitializeDb("db/example.db"); err != nil {
		return err
	}
	return framework.Start(fw)
}
