package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type QuizScoreApiRouter struct {
	qsController *controller.QuizScoreController
	tokenService *auth.TokenService
}

func NewQuizScoreApiRouter(qsController *controller.QuizScoreController, tokenService *auth.TokenService) *QuizScoreApiRouter {
	return &QuizScoreApiRouter{
		qsController: qsController,
		tokenService: tokenService,
	}
}

func (a *QuizScoreApiRouter) Register(r *gin.RouterGroup) {
	fgGroup := r.Group("/final-grades")
	fgGroup.Use(middleware.AuthMiddleware(a.tokenService))
	fgGroup.POST("/:uuid/quiz-scores", a.qsController.CreateQuizScore)
	fgGroup.GET("/:uuid/quiz-scores", a.qsController.ListQuizScores)

	qsGroup := r.Group("/quiz-scores")
	qsGroup.Use(middleware.AuthMiddleware(a.tokenService))
	qsGroup.PUT("/:uuid", a.qsController.UpdateQuizScore)
	qsGroup.DELETE("/:uuid", a.qsController.DeleteQuizScore)
}