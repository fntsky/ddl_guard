package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type DDLApiRouter struct {
	ddlController *controller.DDLController
	tokenService  *auth.TokenService
}

func NewDDLApiRouter(ddlController *controller.DDLController, tokenService *auth.TokenService) *DDLApiRouter {
	return &DDLApiRouter{
		ddlController: ddlController,
		tokenService:  tokenService,
	}
}

func (a *DDLApiRouter) Register(r *gin.RouterGroup) {
	ddlGroup := r.Group("/ddl")
	ddlGroup.Use(middleware.AuthMiddleware(a.tokenService))
	ddlGroup.POST("/draft", a.ddlController.CreateDraft)
	ddlGroup.PATCH("/drafts/:uuid", a.ddlController.ApproveDraft)
}
