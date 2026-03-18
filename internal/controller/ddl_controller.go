package controller

import (
	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/gin-gonic/gin"
)

type DDLController struct {
}

func NewDDLController() *DDLController {
	return &DDLController{}
}

// @Summary 创建DDL草稿
// @Description 创建DDL草稿
// @Tags DDL
// @Accept json
// @Produce json
// @Param req body schema.CreateDraftReq true "Create Draft Request"
// @success 200 {object} handler.resp{data=schema.DraftInfo}
// @Router /ddl/draft [post]
func (dc *DDLController) CreateDraft(ctx *gin.Context) {
	req := &schema.CreateDraftReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
}
