package router

import (
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/google/wire"
)

var ProviderSetRouter = wire.NewSet(
	controller.ProviderSetController,
	NewSwaggerRouter,
	NewAuthApiRouter,
	NewUserApiRouter,
	NewDDLApiRouter,
	NewExamApiRouter,
)
