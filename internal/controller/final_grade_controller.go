package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/final_grade"
	"github.com/gin-gonic/gin"
)

type FinalGradeController struct {
	fgService *final_grade.FinalGradeService
}

func NewFinalGradeController(fgService *final_grade.FinalGradeService) *FinalGradeController {
	return &FinalGradeController{
		fgService: fgService,
	}
}

// @Summary 创建期末成绩
// @Description 创建一个新的期末成绩记录
// @Tags FinalGrade
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body schema.CreateFinalGradeReq true "Create Final Grade Request" SchemaExample({"name":"2024春季期末成绩","exam_ratio":40,"daily_ratio":60})
// @success 200 {object} handler.resp{data=schema.CreateFinalGradeResp} "success"
// @Router /api/v1/final-grades [post]
func (c *FinalGradeController) CreateFinalGrade(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	req := &schema.CreateFinalGradeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := c.fgService.CreateFinalGrade(ctx, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取期末成绩列表
// @Description 分页获取用户所有期末成绩记录
// @Tags FinalGrade
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @success 200 {object} handler.resp{data=schema.FinalGradeListResp} "success"
// @Router /api/v1/final-grades [get]
func (c *FinalGradeController) ListFinalGrades(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	var pageReq schema.PageReq
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid pagination params", err), nil)
		return
	}

	resp, err := c.fgService.ListFinalGrades(ctx, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取期末成绩详情
// @Description 获取单个期末成绩记录的详细信息，包含关联的平时成绩
// @Tags FinalGrade
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Final Grade UUID"
// @success 200 {object} handler.resp{data=schema.FinalGradeDetailResp} "success"
// @Router /api/v1/final-grades/{uuid} [get]
func (c *FinalGradeController) GetFinalGrade(ctx *gin.Context) {
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

	resp, err := c.fgService.GetFinalGrade(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 更新期末成绩
// @Description 更新期末成绩记录的信息
// @Tags FinalGrade
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Final Grade UUID"
// @Param req body schema.UpdateFinalGradeReq true "Update Final Grade Request"
// @success 200 {object} handler.resp{data=schema.UpdateFinalGradeResp} "success"
// @Router /api/v1/final-grades/{uuid} [put]
func (c *FinalGradeController) UpdateFinalGrade(ctx *gin.Context) {
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

	var req schema.UpdateFinalGradeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := c.fgService.UpdateFinalGrade(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除期末成绩
// @Description 删除期末成绩记录及其关联的平时成绩
// @Tags FinalGrade
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Final Grade UUID"
// @success 200 {object} handler.resp "success"
// @Router /api/v1/final-grades/{uuid} [delete]
func (c *FinalGradeController) DeleteFinalGrade(ctx *gin.Context) {
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

	err := c.fgService.DeleteFinalGrade(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}
