package schema

// CreateQuizScoreReq 创建小测成绩请求
type CreateQuizScoreReq struct {
	Name  string  `json:"name" binding:"required" example:"小测1"`
	Score float64 `json:"score" example:"90.0"`
}

// CreateQuizScoreResp 创建小测成绩响应
type CreateQuizScoreResp struct {
	UUID         string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	FinalGradeID string  `json:"final_grade_id" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name         string  `json:"name" example:"小测1"`
	Score        float64 `json:"score" example:"90.0"`
}

// UpdateQuizScoreReq 更新小测成绩请求
type UpdateQuizScoreReq struct {
	Name  *string  `json:"name,omitempty" example:"小测1"`
	Score *float64 `json:"score,omitempty" example:"90.0"`
}

// UpdateQuizScoreResp 更新小测成绩响应
type UpdateQuizScoreResp struct {
	UUID  string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name  string  `json:"name" example:"小测1"`
	Score float64 `json:"score" example:"90.0"`
}

// QuizScoreListItem 小测成绩列表项
type QuizScoreListItem struct {
	UUID     string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name     string  `json:"name" example:"小测1"`
	Score    float64 `json:"score" example:"90.0"`
	CreateAt string  `json:"created_at" example:"2026-04-20T10:00:00+08:00"`
}

// QuizScoreListResp 小测成绩列表响应
type QuizScoreListResp struct {
	List     []QuizScoreListItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}