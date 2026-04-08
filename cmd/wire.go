//go:build wireinject
// +build wireinject

package ddlcmd

import (
	"github.com/fntsky/ddl_guard/internal/base/server"
	"github.com/fntsky/ddl_guard/internal/repo"
	"github.com/fntsky/ddl_guard/internal/router"
	"github.com/fntsky/ddl_guard/internal/service"
	"github.com/fntsky/ddl_guard/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type app struct {
	HttpServer *gin.Engine
	Worker     *worker.PublishWorker
}

func newApp(debug bool, httpServer *gin.Engine, w *worker.PublishWorker) *app {
	return &app{
		HttpServer: httpServer,
		Worker:     w,
	}
}

func (a *app) StartWorker() {
	if a.Worker != nil {
		a.Worker.Start()
	}
}

func initApplication(debug bool) (*app, func(), error) {
	panic(wire.Build(
		service.ProviderSetService,
		repo.ProviderSetRepo,
		server.ProviderSetServer,
		router.ProviderSetRouter,
		worker.NewPublishWorker,
		newApp,
	))
}
