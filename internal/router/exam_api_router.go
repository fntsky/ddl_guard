package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ExamApiRouter struct {
	examController *controller.ExamController
	tokenService   *auth.TokenService
}

func NewExamApiRouter(examController *controller.ExamController, tokenService *auth.TokenService) *ExamApiRouter {
	return &ExamApiRouter{
		examController: examController,
		tokenService:   tokenService,
	}
}

func (a *ExamApiRouter) Register(r *gin.RouterGroup) {
	examGroup := r.Group("/exams")
	examGroup.Use(middleware.AuthMiddleware(a.tokenService))
	examGroup.POST("", a.examController.CreateExam)
	examGroup.GET("", a.examController.ListExams)
	examGroup.GET("/:uuid", a.examController.GetExam)
	examGroup.PUT("/:uuid", a.examController.UpdateExam)
	examGroup.DELETE("/:uuid", a.examController.DeleteExam)
}
