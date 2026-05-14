package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HomeworkScoreApiRouter struct {
	hsController *controller.HomeworkScoreController
	tokenService *auth.TokenService
}

func NewHomeworkScoreApiRouter(hsController *controller.HomeworkScoreController, tokenService *auth.TokenService) *HomeworkScoreApiRouter {
	return &HomeworkScoreApiRouter{
		hsController: hsController,
		tokenService: tokenService,
	}
}

func (a *HomeworkScoreApiRouter) Register(r *gin.RouterGroup) {
	fgGroup := r.Group("/final-grades")
	fgGroup.Use(middleware.AuthMiddleware(a.tokenService))
	fgGroup.POST("/:uuid/homework-scores", a.hsController.CreateHomeworkScore)
	fgGroup.GET("/:uuid/homework-scores", a.hsController.ListHomeworkScores)

	hsGroup := r.Group("/homework-scores")
	hsGroup.Use(middleware.AuthMiddleware(a.tokenService))
	hsGroup.PUT("/:uuid", a.hsController.UpdateHomeworkScore)
	hsGroup.DELETE("/:uuid", a.hsController.DeleteHomeworkScore)
}