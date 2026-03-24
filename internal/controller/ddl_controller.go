package controller

import (
	"strings"

	"github.com/fntsky/ddl_guard/internal/base/handler"
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
// @Param req body schema.CreateDraftReq true "Create Draft Request" SchemaExample({"data_type":"picture","raw_base64":"iVBORw0KGgoAAAANSUhEUgAA...","draft":{"title":"test ddl","description":"from swagger","deadline":"2026-03-24T14:30:00+08:00","early_remind":30}})
// @success 200 {object} handler.resp{data=schema.CreateDraftResp} "success"
// @Router /ddl/draft [post]
func (dc *DDLController) CreateDraft(ctx *gin.Context) {
	req := &schema.CreateDraftReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := dc.ddl_service.CreateDraft(ctx, req)
	handler.HandleResponse(ctx, err, resp)

}

// @Summary 同意DDL草稿
// @Description 将草稿状态从draft变更为active
// @Tags DDL
// @Accept json
// @Produce json
// @Param uuid path string true "Draft UUID"
// @Param req body schema.UpdateDraftStatusReq true "Update Draft Status Request" SchemaExample({"status":"active"})
// @success 200 {object} handler.resp{data=schema.UpdateDraftStatusResp} "success"
// @Router /ddl/drafts/{uuid} [patch]
func (dc *DDLController) ApproveDraft(ctx *gin.Context) {
	uuid := strings.TrimSpace(ctx.Param("uuid"))
	if uuid == "" {
		handler.HandleResponse(ctx, handler.BadRequest("uuid is required", nil), nil)
		return
	}

	req := &schema.UpdateDraftStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := dc.ddl_service.ApproveDraft(ctx, uuid, req)
	handler.HandleResponse(ctx, err, resp)
}
