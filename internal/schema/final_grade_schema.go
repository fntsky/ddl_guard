package schema

// CreateFinalGradeReq 创建期末成绩请求
type CreateFinalGradeReq struct {
	Name       string `json:"name" binding:"required" example:"2024春季期末成绩"`      // 期末成绩名称
	ExamRatio  int    `json:"exam_ratio" example:"40"`                              // 期末考试占比（默认40）
	DailyRatio int    `json:"daily_ratio" example:"60"`                             // 平时成绩占比（默认60）
}

// CreateFinalGradeResp 创建期末成绩响应
type CreateFinalGradeResp struct {
	UUID       string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name       string  `json:"name" example:"2024春季期末成绩"`
	ExamScore  float64 `json:"exam_score" example:"85.5"`
	ExamRatio  int     `json:"exam_ratio" example:"40"`
	DailyRatio int     `json:"daily_ratio" example:"60"`
	FinalScore float64 `json:"final_score" example:"87.0"`
}

// UpdateFinalGradeReq 更新期末成绩请求
type UpdateFinalGradeReq struct {
	Name       *string  `json:"name,omitempty" example:"2024春季期末成绩"`
	ExamScore  *float64 `json:"exam_score,omitempty" example:"85.5"`
	ExamRatio  *int     `json:"exam_ratio,omitempty" example:"40"`
	DailyRatio *int     `json:"daily_ratio,omitempty" example:"60"`
	FinalScore *float64 `json:"final_score,omitempty" example:"87.0"`
}

// UpdateFinalGradeResp 更新期末成绩响应
type UpdateFinalGradeResp struct {
	UUID       string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name       string  `json:"name" example:"2024春季期末成绩"`
	ExamScore  float64 `json:"exam_score" example:"85.5"`
	ExamRatio  int     `json:"exam_ratio" example:"40"`
	DailyRatio int     `json:"daily_ratio" example:"60"`
	FinalScore float64 `json:"final_score" example:"87.0"`
}

// FinalGradeListItem 期末成绩列表项
type FinalGradeListItem struct {
	UUID       string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name       string  `json:"name" example:"2024春季期末成绩"`
	ExamScore  float64 `json:"exam_score" example:"85.5"`
	ExamRatio  int     `json:"exam_ratio" example:"40"`
	DailyRatio int     `json:"daily_ratio" example:"60"`
	FinalScore float64 `json:"final_score" example:"87.0"`
	CreatedAt  string  `json:"created_at" example:"2026-04-20T10:00:00+08:00"`
}

// FinalGradeListResp 期末成绩列表分页响应
type FinalGradeListResp struct {
	List     []FinalGradeListItem `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// FinalGradeDetailResp 期末成绩详情响应
type FinalGradeDetailResp struct {
	UUID       string             `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name       string             `json:"name" example:"2024春季期末成绩"`
	ExamScore  float64            `json:"exam_score" example:"85.5"`
	ExamRatio  int                `json:"exam_ratio" example:"40"`
	DailyRatio int                `json:"daily_ratio" example:"60"`
	FinalScore float64            `json:"final_score" example:"87.0"`
	CreatedAt  string             `json:"created_at" example:"2026-04-20T10:00:00+08:00"`
	UpdatedAt  string             `json:"updated_at" example:"2026-04-20T10:00:00+08:00"`
	DailyScores []DailyScoreItem  `json:"daily_scores,omitempty"` // 关联的平时成绩列表
}

// DailyScoreItem 平时成绩项（嵌入期末成绩详情）
type DailyScoreItem struct {
	UUID  string  `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Type  string  `json:"type" example:"quiz"`        // quiz 或 homework
	Name  string  `json:"name" example:"小测1"`
	Score float64 `json:"score" example:"90.0"`
	Ratio int     `json:"ratio" example:"20"`         // 占平时成绩的比例
}
