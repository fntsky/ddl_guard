package schema

import "time"

const (
	DDLTYPEDEFAULT = "default"
	DDLTYPEPICTURE = "picture"
)

type CreateDraftReq struct {
	Raw   []byte          `json:"raw"`
	Type  string          `json:"data_type"`
	Draft CreateDraftResp `json:"draft"`
}

/*
uuid
标题
描述
截止时间
截止时间前多少分钟提醒
*/
type CreateDraftResp struct {
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	EealyRemind int       `json:"early_remind"`
}
