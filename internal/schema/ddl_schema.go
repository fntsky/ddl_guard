package schema

import "time"

const (
	DDLTYPEDEFAULT  = "default"
	DDLTYPEPICTURE  = "picture"
	DDLSTATUSACTIVE = "active"
)

type CreateDraftReq struct {
	RawBase64 string           `json:"raw_base64,omitempty" example:"iVBORw0KGgoAAAANSUhEUgAA..."`
	Type      string           `json:"data_type" default:"default" example:"default"`
	Draft     CreateDraftInput `json:"draft"`
}

type CreateDraftInput struct {
	Title       string    `json:"title" example:"test ddl"`
	Description string    `json:"description" example:"from swagger"`
	Deadline    time.Time `json:"deadline" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemind int       `json:"early_remind" example:"30"`
}

/*
uuid
标题
描述
截止时间
截止时间前多少分钟提醒
*/
type CreateDraftResp struct {
	UUID        string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Title       string    `json:"title" example:"test ddl"`
	Description string    `json:"description" example:"from swagger"`
	Deadline    time.Time `json:"deadline" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemind int       `json:"early_remind" example:"30"`
}

type UpdateDraftStatusReq struct {
	Status string `json:"status" binding:"required" example:"active"`
}

type UpdateDraftStatusResp struct {
	UUID   string `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Status int    `json:"status" example:"1"`
}

// PageReq 分页请求
type PageReq struct {
	Page     int `form:"page" example:"1"`
	PageSize int `form:"page_size" example:"10"`
}

// Normalize 设置默认值
func (r *PageReq) Normalize() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 {
		r.PageSize = 10
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
}

// Offset 计算偏移量
func (r *PageReq) Offset() int {
	return (r.Page - 1) * r.PageSize
}

// DDLListItem DDL列表项
type DDLListItem struct {
	UUID            string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Title           string    `json:"title" example:"test ddl"`
	Description     string    `json:"description" example:"from swagger"`
	Deadline        time.Time `json:"deadline" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemindTime time.Time `json:"early_remind_time" swaggertype:"string" format:"date-time" example:"2026-03-24T14:00:00+08:00"`
	Status          int       `json:"status" example:"1"`
	CreatedAt       time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2026-03-24T10:00:00+08:00"`
}

// DDLListResp DDL列表分页响应
type DDLListResp struct {
	List     []DDLListItem `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// UpdateDDLReq 修改DDL请求
type UpdateDDLReq struct {
	Title       *string    `json:"title,omitempty" example:"test ddl"`
	Description *string    `json:"description,omitempty" example:"from swagger"`
	Deadline    *time.Time `json:"deadline,omitempty" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemind *int       `json:"early_remind,omitempty" example:"30"` // 提前多少分钟提醒
}

// UpdateDDLResp 修改DDL响应
type UpdateDDLResp struct {
	UUID        string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Title       string    `json:"title" example:"test ddl"`
	Description string    `json:"description" example:"from swagger"`
	Deadline    time.Time `json:"deadline" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemind int       `json:"early_remind" example:"30"`
}

// DDLDetailResp DDL详情响应
type DDLDetailResp struct {
	UUID            string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Title           string    `json:"title" example:"test ddl"`
	Description     string    `json:"description" example:"from swagger"`
	Deadline        time.Time `json:"deadline" swaggertype:"string" format:"date-time" example:"2026-03-24T14:30:00+08:00"`
	EarlyRemindTime time.Time `json:"early_remind_time" swaggertype:"string" format:"date-time" example:"2026-03-24T14:00:00+08:00"`
	Status          int       `json:"status" example:"1"`
	RemindSent      bool      `json:"remind_sent" example:"false"`
	CreatedAt       time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2026-03-24T10:00:00+08:00"`
	UpdatedAt       time.Time `json:"updated_at" swaggertype:"string" format:"date-time" example:"2026-03-24T10:00:00+08:00"`
}
