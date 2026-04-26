package entity

import "time"

// DailyScoreType 平时成绩类型
type DailyScoreType string

const (
	DailyScoreTypeQuiz     DailyScoreType = "quiz"     // 平时小测
	DailyScoreTypeHomework DailyScoreType = "homework" // 平时作业
)

// DailyScore 平时成绩
type DailyScore struct {
	ID           int64          `xorm:"pk autoincr 'id'"`
	UUID         string         `xorm:"uuid not null unique index 'uuid'"`
	FinalGradeID int64          `xorm:"not null index 'final_grade_id'"` // 关联的期末成绩记录
	UserID       int64          `xorm:"not null index 'user_id'"`
	Type         DailyScoreType `xorm:"varchar(20) not null 'type'"`     // 类型：quiz/homework
	Name         string         `xorm:"text not null 'name'"`             // 名称，如"小测1"、"作业1"
	Score        float64        `xorm:"default 0 'score'"`                // 得分（满分100）
	Ratio        int            `xorm:"default 0 'ratio'"`                // 占平时成绩的比例（百分比）
	CreatedAt    time.Time      `xorm:"created 'created_at'"`
	UpdatedAt    time.Time      `xorm:"updated 'updated_at'"`
}

func (DailyScore) TableName() string {
	return "daily_scores"
}