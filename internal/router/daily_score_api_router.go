package router

import (
	"github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/gin-gonic/gin"
)

type DailyScoreApiRouter struct {
	dsController *controller.DailyScoreController
	tokenService *auth.TokenService
}

func NewDailyScoreApiRouter(dsController *controller.DailyScoreController, tokenService *auth.TokenService) *DailyScoreApiRouter {
	return &DailyScoreApiRouter{
		dsController: dsController,
		tokenService: tokenService,
	}
}

func (a *DailyScoreApiRouter) Register(r *gin.RouterGroup) {
	// 平时成绩嵌套在期末成绩路由下
	fgGroup := r.Group("/final-grades")
	fgGroup.Use(middleware.AuthMiddleware(a.tokenService))
	fgGroup.POST("/:final_grade_uuid/daily-scores", a.dsController.CreateDailyScore)
	fgGroup.GET("/:final_grade_uuid/daily-scores", a.dsController.ListDailyScores)

	// 平时成绩的更新和删除使用独立路由
	dsGroup := r.Group("/daily-scores")
	dsGroup.Use(middleware.AuthMiddleware(a.tokenService))
	dsGroup.PUT("/:uuid", a.dsController.UpdateDailyScore)
	dsGroup.DELETE("/:uuid", a.dsController.DeleteDailyScore)
}