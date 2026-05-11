package server

import (
	"regexp"

	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
			field := fl.Field().String()
			return mobileRegex.MatchString(field)
		})
	}
}

func NewHttpServer(debug bool,
	swaggerRouter *router.SwaggerRouter,
	authApiRouter *router.AuthApiRouter,
	ddlApiRouter *router.DDLApiRouter,
	examApiRouter *router.ExamApiRouter,
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
	examApiRouter.Register(apiGroup)
	return r
}
