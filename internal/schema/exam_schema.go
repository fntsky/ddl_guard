package schema

import "time"

// CreateExamReq 创建考试请求
type CreateExamReq struct {
	Name      string    `json:"name" binding:"required" example:"高等数学期末考试"`
	StartTime time.Time `json:"start_time" binding:"required" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   time.Time `json:"end_time" binding:"required" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  string    `json:"location" example:"教学楼A301"`
	Notes     string    `json:"notes" example:"带计算器"`
}

// CreateExamResp 创建考试响应
type CreateExamResp struct {
	UUID      string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name      string    `json:"name" example:"高等数学期末考试"`
	StartTime time.Time `json:"start_time" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   time.Time `json:"end_time" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  string    `json:"location" example:"教学楼A301"`
	Notes     string    `json:"notes" example:"带计算器"`
}

// UpdateExamReq 修改考试请求
type UpdateExamReq struct {
	Name      *string    `json:"name,omitempty" example:"高等数学期末考试"`
	StartTime *time.Time `json:"start_time,omitempty" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   *time.Time `json:"end_time,omitempty" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  *string    `json:"location,omitempty" example:"教学楼A301"`
	Notes     *string    `json:"notes,omitempty" example:"带计算器"`
}

// UpdateExamResp 修改考试响应
type UpdateExamResp struct {
	UUID      string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name      string    `json:"name" example:"高等数学期末考试"`
	StartTime time.Time `json:"start_time" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   time.Time `json:"end_time" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  string    `json:"location" example:"教学楼A301"`
	Notes     string    `json:"notes" example:"带计算器"`
}

// ExamListItem 考试列表项
type ExamListItem struct {
	UUID      string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name      string    `json:"name" example:"高等数学期末考试"`
	StartTime time.Time `json:"start_time" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   time.Time `json:"end_time" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  string    `json:"location" example:"教学楼A301"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2026-04-20T10:00:00+08:00"`
}

// ExamListResp 考试列表分页响应
type ExamListResp struct {
	List     []ExamListItem `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// ExamDetailResp 考试详情响应
type ExamDetailResp struct {
	UUID      string    `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Name      string    `json:"name" example:"高等数学期末考试"`
	StartTime time.Time `json:"start_time" swaggertype:"string" format:"date-time" example:"2026-04-26T09:00:00+08:00"`
	EndTime   time.Time `json:"end_time" swaggertype:"string" format:"date-time" example:"2026-04-26T11:00:00+08:00"`
	Location  string    `json:"location" example:"教学楼A301"`
	Notes     string    `json:"notes" example:"带计算器"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2026-04-20T10:00:00+08:00"`
	UpdatedAt time.Time `json:"updated_at" swaggertype:"string" format:"date-time" example:"2026-04-20T10:00:00+08:00"`
}
