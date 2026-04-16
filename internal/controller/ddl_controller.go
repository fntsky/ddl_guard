package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/middleware"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	"github.com/gin-gonic/gin"
)

type DDLController struct {
	ddl_service *ddl.DDLService
}

func NewDDLController(ddl_service *ddl.DDLService) *DDLController {
	return &DDLController{
		ddl_service: ddl_service,
	}
}

// @Summary 创建DDL草稿
// @Description 创建DDL草稿
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body schema.CreateDraftReq true "Create Draft Request" SchemaExample({"data_type":"picture","raw_base64":"iVBORw0KGgoAAAANSUhEUgAA...","draft":{"title":"test ddl","description":"from swagger","deadline":"2026-03-24T14:30:00+08:00","early_remind":30}})
// @success 200 {object} handler.resp{data=schema.CreateDraftResp} "success"
// @Router /api/v1/ddl/draft [post]
func (dc *DDLController) CreateDraft(ctx *gin.Context) {
	userClaims, ok := middleware.GetUserFromGin(ctx)
	if !ok || userClaims.UserUUID == "" {
		handler.HandleResponse(ctx, handler.NewError(401, "unauthorized", nil), nil)
		return
	}

	req := &schema.CreateDraftReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := dc.ddl_service.CreateDraft(ctx, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)

}

// @Summary 同意DDL草稿
// @Description 将草稿状态从draft变更为active
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Draft UUID"
// @Param req body schema.UpdateDraftStatusReq true "Update Draft Status Request" SchemaExample({"status":"active"})
// @success 200 {object} handler.resp{data=schema.UpdateDraftStatusResp} "success"
// @Router /api/v1/ddl/drafts/{uuid} [patch]
func (dc *DDLController) ApproveDraft(ctx *gin.Context) {
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

	req := &schema.UpdateDraftStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := dc.ddl_service.ApproveDraft(ctx, uuid, req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 删除DDL
// @Description 删除DDL（软删除）
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "DDL UUID"
// @success 200 {object} handler.resp "success"
// @Router /api/v1/ddl/{uuid} [delete]
func (dc *DDLController) DeleteDDL(ctx *gin.Context) {
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

	err := dc.ddl_service.DeleteDDL(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, nil)
}

// @Summary 获取激活状态的DDL列表
// @Description 分页获取用户所有激活状态的DDL
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @success 200 {object} handler.resp{data=schema.DDLListResp} "success"
// @Router /api/v1/ddl/active [get]
func (dc *DDLController) GetActiveDDLs(ctx *gin.Context) {
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

	resp, err := dc.ddl_service.GetActiveDDLs(ctx, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取过期状态的DDL列表
// @Description 分页获取用户所有过期状态的DDL
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @success 200 {object} handler.resp{data=schema.DDLListResp} "success"
// @Router /api/v1/ddl/expired [get]
func (dc *DDLController) GetExpiredDDLs(ctx *gin.Context) {
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

	resp, err := dc.ddl_service.GetExpiredDDLs(ctx, userClaims.UserUUID, &pageReq)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 修改DDL
// @Description 修改DDL的标题、描述、截止时间或提前提醒时间
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "DDL UUID"
// @Param req body schema.UpdateDDLReq true "Update DDL Request"
// @success 200 {object} handler.resp{data=schema.UpdateDDLResp} "success"
// @Router /api/v1/ddl/{uuid} [put]
func (dc *DDLController) UpdateDDL(ctx *gin.Context) {
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

	var req schema.UpdateDDLReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(ctx, handler.BadRequest("invalid request body", err), nil)
		return
	}

	resp, err := dc.ddl_service.UpdateDDL(ctx, uuid, &req, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 获取DDL详情
// @Description 获取单个DDL的详细信息
// @Tags DDL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "DDL UUID"
// @success 200 {object} handler.resp{data=schema.DDLDetailResp} "success"
// @Router /api/v1/ddl/{uuid} [get]
func (dc *DDLController) GetDDLDetail(ctx *gin.Context) {
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

	resp, err := dc.ddl_service.GetDDLDetail(ctx, uuid, userClaims.UserUUID)
	handler.HandleResponse(ctx, err, resp)
}
