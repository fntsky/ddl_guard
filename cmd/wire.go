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
	HttpServer       *gin.Engine
	PublishWorker    *worker.PublishWorker
	ExpirationWorker *worker.ExpirationWorker
}

func newApp(debug bool, httpServer *gin.Engine, pw *worker.PublishWorker, ew *worker.ExpirationWorker) *app {
	return &app{
		HttpServer:       httpServer,
		PublishWorker:    pw,
		ExpirationWorker: ew,
	}
}

func (a *app) StartWorkers() {
	if a.PublishWorker != nil {
		a.PublishWorker.Start()
	}
	if a.ExpirationWorker != nil {
		a.ExpirationWorker.Start()
	}
}

func initApplication(debug bool) (*app, func(), error) {
	panic(wire.Build(
		service.ProviderSetService,
		repo.ProviderSetRepo,
		server.ProviderSetServer,
		router.ProviderSetRouter,
		worker.NewPublishWorker,
		worker.NewExpirationWorker,
		newApp,
	))
}
