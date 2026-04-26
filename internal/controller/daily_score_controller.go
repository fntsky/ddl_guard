package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/daily_score"
	"github.com/gin-gonic/gin"
)

type DailyScoreController struct {
	dsService *daily_score.DailyScoreService
}

func NewDailyScoreController(dsService *daily_score.DailyScoreService) *DailyScoreController {
	return &DailyScoreController{
		dsService: dsService,
	}
}

// @Summary 创建平时成绩
// @Description 为指定期末成绩记录添加平时成绩（小测或作业）
// @Tags DailyScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @Param req body schema.CreateDailyScoreReq true "Create Daily Score Request" SchemaExample({"type":"quiz","name":"小测1","score":90,"ratio":20})
// @success 200 {object} handler.resp{data=schema.CreateDailyScoreResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/daily-scores [post]
func (c *DailyScoreController) CreateDailyScore(ctx *gin.Context) {
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

	req := &schema.CreateDailyScoreReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := c.dsService.CreateDailyScore(ctx, fgUUID, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取平时成绩列表
// @Description 获取指定期末成绩记录下的所有平时成绩
// @Tags DailyScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param final_grade_uuid path string true "Final Grade UUID"
// @success 200 {object} handler.resp{data=schema.DailyScoreListResp} "success"
// @Router /api/v1/final-grades/{final_grade_uuid}/daily-scores [get]
func (c *DailyScoreController) ListDailyScores(ctx *gin.Context) {
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

	resp, err := c.dsService.ListDailyScores(ctx, fgUUID, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 更新平时成绩
// @Description 更新平时成绩的信息
// @Tags DailyScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Daily Score UUID"
// @Param req body schema.UpdateDailyScoreReq true "Update Daily Score Request"
// @success 200 {object} handler.resp{data=schema.UpdateDailyScoreResp} "success"
// @Router /api/v1/daily-scores/{uuid} [put]
func (c *DailyScoreController) UpdateDailyScore(ctx *gin.Context) {
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

	var req schema.UpdateDailyScoreReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := c.dsService.UpdateDailyScore(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除平时成绩
// @Description 删除平时成绩
// @Tags DailyScore
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Daily Score UUID"
// @success 200 {object} handler.resp "success"
// @Router /api/v1/daily-scores/{uuid} [delete]
func (c *DailyScoreController) DeleteDailyScore(ctx *gin.Context) {
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

	err := c.dsService.DeleteDailyScore(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}
