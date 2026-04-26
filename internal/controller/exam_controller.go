package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/exam"
	"github.com/gin-gonic/gin"
)

type ExamController struct {
	exam_service *exam.ExamService
}

func NewExamController(exam_service *exam.ExamService) *ExamController {
	return &ExamController{
		exam_service: exam_service,
	}
}

// @Summary 创建考试
// @Description 创建一个新的考试
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body schema.CreateExamReq true "Create Exam Request" SchemaExample({"name":"高等数学期末考试","start_time":"2026-04-26T09:00:00+08:00","end_time":"2026-04-26T11:00:00+08:00","location":"教学楼A301","notes":"带计算器"})
// @success 200 {object} handler.resp{data=schema.CreateExamResp} "success"
// @Router /api/v1/exams [post]
func (ec *ExamController) CreateExam(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	req := &schema.CreateExamReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := ec.exam_service.CreateExam(ctx, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取考试列表
// @Description 分页获取用户所有考试
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @success 200 {object} handler.resp{data=schema.ExamListResp} "success"
// @Router /api/v1/exams [get]
func (ec *ExamController) ListExams(ctx *gin.Context) {
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

	resp, err := ec.exam_service.ListExams(ctx, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取考试详情
// @Description 获取单个考试的详细信息
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Exam UUID"
// @success 200 {object} handler.resp{data=schema.ExamDetailResp} "success"
// @Router /api/v1/exams/{uuid} [get]
func (ec *ExamController) GetExam(ctx *gin.Context) {
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

	resp, err := ec.exam_service.GetExam(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 更新考试
// @Description 更新考试的信息
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Exam UUID"
// @Param req body schema.UpdateExamReq true "Update Exam Request"
// @success 200 {object} handler.resp{data=schema.UpdateExamResp} "success"
// @Router /api/v1/exams/{uuid} [put]
func (ec *ExamController) UpdateExam(ctx *gin.Context) {
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

	var req schema.UpdateExamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := ec.exam_service.UpdateExam(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除考试
// @Description 删除考试
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Exam UUID"
// @success 200 {object} handler.resp "success"
// @Router /api/v1/exams/{uuid} [delete]
func (ec *ExamController) DeleteExam(ctx *gin.Context) {
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

	err := ec.exam_service.DeleteExam(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}
