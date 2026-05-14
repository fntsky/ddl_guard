package schema

// CreateHomeworkScoreReq 创建作业成绩请求
type CreateHomeworkScoreReq struct {
	Name  string  `json:"name" binding:"required" example:"作业1"`
	Score float64 `json:"score" example:"90.0"`
}

// CreateHomeworkScoreResp 创建作业成绩响应
type CreateHomeworkScoreResp struct {
	UUID         string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	FinalGradeID string  `json:"final_grade_id" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name         string  `json:"name" example:"作业1"`
	Score        float64 `json:"score" example:"90.0"`
}

// UpdateHomeworkScoreReq 更新作业成绩请求
type UpdateHomeworkScoreReq struct {
	Name  *string  `json:"name,omitempty" example:"作业1"`
	Score *float64 `json:"score,omitempty" example:"90.0"`
}

// UpdateHomeworkScoreResp 更新作业成绩响应
type UpdateHomeworkScoreResp struct {
	UUID  string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name  string  `json:"name" example:"作业1"`
	Score float64 `json:"score" example:"90.0"`
}

// HomeworkScoreListItem 作业成绩列表项
type HomeworkScoreListItem struct {
	UUID     string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name     string  `json:"name" example:"作业1"`
	Score    float64 `json:"score" example:"90.0"`
	CreateAt string  `json:"created_at" example:"2026-04-20T10:00:00+08:00"`
}

// HomeworkScoreListResp 作业成绩列表响应
type HomeworkScoreListResp struct {
	List     []HomeworkScoreListItem `json:"list"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}