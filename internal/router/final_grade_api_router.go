package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type FinalGradeApiRouter struct {
	fgController *controller.FinalGradeController
	tokenService *auth.TokenService
}

func NewFinalGradeApiRouter(fgController *controller.FinalGradeController, tokenService *auth.TokenService) *FinalGradeApiRouter {
	return &FinalGradeApiRouter{
		fgController: fgController,
		tokenService: tokenService,
	}
}

func (a *FinalGradeApiRouter) Register(r *gin.RouterGroup) {
	fgGroup := r.Group("/final-grades")
	fgGroup.Use(middleware.AuthMiddleware(a.tokenService))
	fgGroup.POST("", a.fgController.CreateFinalGrade)
	fgGroup.GET("", a.fgController.ListFinalGrades)
	fgGroup.GET("/:uuid", a.fgController.GetFinalGrade)
	fgGroup.PUT("/:uuid", a.fgController.UpdateFinalGrade)
	fgGroup.DELETE("/:uuid", a.fgController.DeleteFinalGrade)
}