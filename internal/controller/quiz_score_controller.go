package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/quiz_score"
	"github.com/gin-gonic/gin"
)

type QuizScoreController struct {
	qsService *quiz_score.QuizScoreService
}

func NewQuizScoreController(qsService *quiz_score.QuizScoreService) *QuizScoreController {
	return &QuizScoreController{
		qsService: qsService,
	}
}

// @Summary 创建小测成绩
// @Description 为指定期末成绩记录添加小测成绩
// @Tags QuizScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @Param req body schema.CreateQuizScoreReq true "Create Quiz Score Request" SchemaExample({"name":"小测1","score":90})
// @success 200 {object} handler.Response{data=schema.CreateQuizScoreResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/quiz-scores [post]
func (c *QuizScoreController) CreateQuizScore(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	fgUUID := strings.TrimSpace(ctx.Param("uuid"))
	if fgUUID == "" {
		handler.HandleResponse(ctx, handler.BadRequest("uuid is required", nil), nil)
		return
	}

	req := &schema.CreateQuizScoreReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := c.qsService.CreateQuizScore(ctx, fgUUID, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取小测成绩列表
// @Description 获取指定期末成绩记录下的所有小测成绩
// @Tags QuizScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @success 200 {object} handler.Response{data=schema.QuizScoreListResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/quiz-scores [get]
func (c *QuizScoreController) ListQuizScores(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	fgUUID := strings.TrimSpace(ctx.Param("uuid"))
	if fgUUID == "" {
		handler.HandleResponse(ctx, handler.BadRequest("uuid is required", nil), nil)
		return
	}

	var pageReq schema.PageReq
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid pagination params", err), nil)
		return
	}
	pageReq.Normalize()

	resp, err := c.qsService.ListQuizScores(ctx, fgUUID, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 更新小测成绩
// @Description 更新小测成绩的信息，自动重算最终成绩
// @Tags QuizScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Quiz Score UUID"
// @Param req body schema.UpdateQuizScoreReq true "Update Quiz Score Request"
// @success 200 {object} handler.Response{data=schema.UpdateQuizScoreResp} "success"
// @Router /api/v1/quiz-scores/{uuid} [put]
func (c *QuizScoreController) UpdateQuizScore(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	uuid := strings.TrimSpace(ctx.Param("uuid"))
	if uuid == "" {
		handler.HandleResponse(ctx, handler.BadRequest("uuid is required", nil), nil)
		return
	}

	var req schema.UpdateQuizScoreReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := c.qsService.UpdateQuizScore(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除小测成绩
// @Description 删除小测成绩，自动重算最终成绩
// @Tags QuizScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Quiz Score UUID"
// @success 200 {object} handler.Response "success"
// @Router /api/v1/quiz-scores/{uuid} [delete]
func (c *QuizScoreController) DeleteQuizScore(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	uuid := strings.TrimSpace(ctx.Param("uuid"))
	if uuid == "" {
		handler.HandleResponse(ctx, handler.BadRequest("uuid is required", nil), nil)
		return
	}

	err := c.qsService.DeleteQuizScore(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}