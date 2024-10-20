package router

import (
	"net/http"

	"gorm-gen-skeleton/app/controller"
	"gorm-gen-skeleton/app/middleware"

	"gorm-gen-skeleton/internal/server"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	server server.HttpServer
}

func New(server server.HttpServer) *AppRouter {
	server.SetMiddleware(&middleware.Foo{}, &middleware.Cors{})
	return &AppRouter{
		server,
	}
}

func (*AppRouter) Add(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello word!")
	})

	index := &controller.Index{}
	server.GET("/hello", index.Hello)

	server.GET("/socket", (&controller.Socket{}).Connect)
}
