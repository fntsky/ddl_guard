package entity

import "time"

// FinalGrade 期末成绩记录
type FinalGrade struct {
	ID         int64     `xorm:"pk autoincr 'id'"`
	UUID       string    `xorm:"uuid not null unique index 'uuid'"`
	UserID     int64     `xorm:"not null index 'user_id'"`
	Name       string    `xorm:"text not null 'name'"`           // 期末成绩名称，如"2024春季期末成绩"
	ExamScore  float64   `xorm:"default 0 'exam_score'"`         // 期末考试成绩（满分100）
	ExamRatio  int       `xorm:"default 40 'exam_ratio'"`        // 期末考试成绩占比（百分比）
	DailyRatio int       `xorm:"default 60 'daily_ratio'"`       // 平时成绩占比（百分比）
	FinalScore float64   `xorm:"default 0 'final_score'"`        // 最终成绩（用户手动录入）
	CreatedAt  time.Time `xorm:"created index 'created_at'"`
	UpdatedAt  time.Time `xorm:"updated 'updated_at'"`
}

func (FinalGrade) TableName() string {
	return "final_grades"
}