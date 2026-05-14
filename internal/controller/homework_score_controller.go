package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/homework_score"
	"github.com/gin-gonic/gin"
)

type HomeworkScoreController struct {
	hsService *homework_score.HomeworkScoreService
}

func NewHomeworkScoreController(hsService *homework_score.HomeworkScoreService) *HomeworkScoreController {
	return &HomeworkScoreController{
		hsService: hsService,
	}
}

// @Summary 创建作业成绩
// @Description 为指定期末成绩记录添加作业成绩
// @Tags HomeworkScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @Param req body schema.CreateHomeworkScoreReq true "Create Homework Score Request" SchemaExample({"name":"作业1","score":90})
// @success 200 {object} handler.Response{data=schema.CreateHomeworkScoreResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/homework-scores [post]
func (c *HomeworkScoreController) CreateHomeworkScore(ctx *gin.Context) {
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

	req := &schema.CreateHomeworkScoreReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := c.hsService.CreateHomeworkScore(ctx, fgUUID, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取作业成绩列表
// @Description 获取指定期末成绩记录下的所有作业成绩
// @Tags HomeworkScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @success 200 {object} handler.Response{data=schema.HomeworkScoreListResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/homework-scores [get]
func (c *HomeworkScoreController) ListHomeworkScores(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	fgUUID := strings.TrimSpace(ctx.Param("final_grade_uuid"))
	if fgUUID == "" {
		handler.HandleResponse(ctx, handler.BadRequest("final_grade_uuid is required", nil), nil)
		return
	}

	var pageReq schema.PageReq
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid pagination params", err), nil)
		return
	}
	pageReq.Normalize()

	resp, err := c.hsService.ListHomeworkScores(ctx, fgUUID, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 更新作业成绩
// @Description 更新作业成绩的信息，自动重算最终成绩
// @Tags HomeworkScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Homework Score UUID"
// @Param req body schema.UpdateHomeworkScoreReq true "Update Homework Score Request"
// @success 200 {object} handler.Response{data=schema.UpdateHomeworkScoreResp} "success"
// @Router /api/v1/homework-scores/{uuid} [put]
func (c *HomeworkScoreController) UpdateHomeworkScore(ctx *gin.Context) {
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

	var req schema.UpdateHomeworkScoreReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := c.hsService.UpdateHomeworkScore(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除作业成绩
// @Description 删除作业成绩，自动重算最终成绩
// @Tags HomeworkScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Homework Score UUID"
// @success 200 {object} handler.Response "success"
// @Router /api/v1/homework-scores/{uuid} [delete]
func (c *HomeworkScoreController) DeleteHomeworkScore(ctx *gin.Context) {
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

	err := c.hsService.DeleteHomeworkScore(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}