package router

import (
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/gin-gonic/gin"
)

type DDLApiRouter struct {
	ddlController *controller.DDLController
}

func NewDDLApiRouter(ddlController *controller.DDLController) *DDLApiRouter {
	return &DDLApiRouter{ddlController: ddlController}
}

func (a *DDLApiRouter) Register(r *gin.RouterGroup) {
	ddlGroup := r.Group("/ddl")
	ddlGroup.POST("/draft", a.ddlController.CreateDraft)
}
