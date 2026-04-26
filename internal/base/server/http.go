package server

import (
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/router"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(debug bool,
	swaggerRouter *router.SwaggerRouter,
	authApiRouter *router.AuthApiRouter,
	ddlApiRouter *router.DDLApiRouter,
	userApiRouter *router.UserApiRouter) *gin.Engine {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 注册请求 ID 中间件（放在最前面）
	r.Use(middleware.RequestIDMiddleware())

	rootGroup := r.Group("")
	apiGroup := r.Group("/api/v1")
	swaggerRouter.Register(rootGroup)
	authApiRouter.Register(apiGroup)
	userApiRouter.Register(apiGroup)
	ddlApiRouter.Register(apiGroup)
	return r
}
