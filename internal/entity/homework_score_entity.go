package entity

import "time"

// HomeworkScore 作业成绩
type HomeworkScore struct {
	ID           int64     `xorm:"pk autoincr 'id'"`
	UUID         string    `xorm:"uuid not null unique index 'uuid'"`
	FinalGradeID int64     `xorm:"not null index 'final_grade_id'"`
	UserID       int64     `xorm:"not null index 'user_id'"`
	Name         string    `xorm:"text not null 'name'"`             // 名称，如"作业1"
	Score        float64   `xorm:"default 0 'score'"`                // 得分（满分100）
	CreatedAt    time.Time `xorm:"created 'created_at'"`
	UpdatedAt    time.Time `xorm:"updated 'updated_at'"`
}

func (HomeworkScore) TableName() string {
	return "homework_scores"
}
