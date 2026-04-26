package schema

// CreateDailyScoreReq 创建平时成绩请求
type CreateDailyScoreReq struct {
	Type  string  `json:"type" binding:"required,oneof=quiz homework" example:"quiz"` // quiz 或 homework
	Name  string  `json:"name" binding:"required" example:"小测1"`
	Score float64 `json:"score" example:"90.0"`
	Ratio int     `json:"ratio" example:"20"` // 占平时成绩的比例
}

// CreateDailyScoreResp 创建平时成绩响应
type CreateDailyScoreResp struct {
	UUID         string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	FinalGradeID string  `json:"final_grade_id" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Type         string  `json:"type" example:"quiz"`
	Name         string  `json:"name" example:"小测1"`
	Score        float64 `json:"score" example:"90.0"`
	Ratio        int     `json:"ratio" example:"20"`
}

// UpdateDailyScoreReq 更新平时成绩请求
type UpdateDailyScoreReq struct {
	Type  *string  `json:"type,omitempty" example:"quiz"`
	Name  *string  `json:"name,omitempty" example:"小测1"`
	Score *float64 `json:"score,omitempty" example:"90.0"`
	Ratio *int     `json:"ratio,omitempty" example:"20"`
}

// UpdateDailyScoreResp 更新平时成绩响应
type UpdateDailyScoreResp struct {
	UUID  string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Type  string  `json:"type" example:"quiz"`
	Name  string  `json:"name" example:"小测1"`
	Score float64 `json:"score" example:"90.0"`
	Ratio int     `json:"ratio" example:"20"`
}

// DailyScoreListItem 平时成绩列表项
type DailyScoreListItem struct {
	UUID     string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Type     string  `json:"type" example:"quiz"`
	Name     string  `json:"name" example:"小测1"`
	Score    float64 `json:"score" example:"90.0"`
	Ratio    int     `json:"ratio" example:"20"`
	CreateAt string  `json:"created_at" example:"2026-04-20T10:00:00+08:00"`
}

// DailyScoreListResp 平时成绩列表响应
type DailyScoreListResp struct {
	List     []DailyScoreListItem `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}
