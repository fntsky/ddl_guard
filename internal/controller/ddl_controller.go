package controller

import (
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
// @Param req body schema.CreateDraftReq true "Create Draft Request"
// @success 200 {object} handler.resp{data=schema.CreateDraftResp} "success"
// @Router /ddl/draft [post]
func (dc *DDLController) CreateDraft(ctx *gin.Context) {
	req := &schema.CreateDraftReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := dc.ddl_service.CreateDraft(ctx, req)
	handler.HanderResponse(ctx, err, resp)

}
