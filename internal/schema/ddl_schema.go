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
	EealyRemind int       `json:"early_remind" example:"30"`
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
	EealyRemind int       `json:"early_remind" example:"30"`
}

type UpdateDraftStatusReq struct {
	Status string `json:"status" binding:"required" example:"active"`
}

type UpdateDraftStatusResp struct {
	UUID   string `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Status int    `json:"status" example:"1"`
}
