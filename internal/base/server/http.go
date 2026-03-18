package server

import (
	"github.com/fntsky/ddl_guard/internal/router"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(debug bool,
	swaggerRouter *router.SwaggerRouter,
	ddlApiRouter *router.DDLApiRouter) *gin.Engine {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	rootGroup := r.Group("")
	apiGroup := r.Group("/api/v1")
	swaggerRouter.Register(rootGroup)
	ddlApiRouter.Register(apiGroup)
	return r
}
