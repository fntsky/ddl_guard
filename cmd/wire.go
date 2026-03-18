//go:build wireinject
// +build wireinject

package ddlcmd

import (
	"github.com/fntsky/ddl_guard/internal/base/server"
	"github.com/fntsky/ddl_guard/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type app struct {
	HttpServer *gin.Engine
}

func newApp(debug bool, httpServer *gin.Engine) *app {
	return &app{
		HttpServer: httpServer,
	}
}

func initApplication(debug bool) (*app, func(), error) {
	panic(wire.Build(
		server.ProviderSetServer,
		router.ProviderSetRouter,
		newApp,
	))
}
