package main

import (
	_ "gorm-gen-skeleton/internal/bootstrap"
	"gorm-gen-skeleton/internal/server"
	"gorm-gen-skeleton/internal/variable"
	"gorm-gen-skeleton/router"
)

func main() {
	variable.Init()
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")
	http := server.New(
		server.WithPort(port),
		server.WithMode(mode),
		server.WithLogger(variable.Log),
	)
	http.SetRouters(router.New(http)).Run()
}
