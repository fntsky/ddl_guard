package entity

import "time"

// FinalGrade 期末成绩记录
type FinalGrade struct {
	ID                    int64     `xorm:"pk autoincr 'id'"`
	UUID                  string    `xorm:"uuid not null unique index 'uuid'"`
	UserID                int64     `xorm:"not null index 'user_id'"`
	Name                  string    `xorm:"text not null 'name'"`                            // 期末成绩名称
	ExamScore             float64   `xorm:"default 0 'exam_score'"`                          // 期末考试成绩（满分100）
	ExamRatio             int       `xorm:"default 0 'exam_ratio'"`                          // 期末考试占比
	ClassroomBonusScore   float64   `xorm:"default 0 'classroom_bonus_score'"`               // 课堂加分（满分100）
	ClassroomBonusRatio   int       `xorm:"default 0 'classroom_bonus_ratio'"`               // 课堂加分占比
	AttendanceScore       float64   `xorm:"default 0 'attendance_score'"`                     // 考勤成绩（满分100）
	AttendanceRatio       int       `xorm:"default 0 'attendance_ratio'"`                     // 考勤占比
	QuizRatio             int       `xorm:"default 0 'quiz_ratio'"`                           // 小测占比
	HomeworkRatio         int       `xorm:"default 0 'homework_ratio'"`                       // 作业占比
	FinalScore            float64   `xorm:"default 0 'final_score'"`                          // 最终成绩（自动计算）
	CreatedAt             time.Time `xorm:"created index 'created_at'"`
	UpdatedAt             time.Time `xorm:"updated 'updated_at'"`
}

func (FinalGrade) TableName() string {
	return "final_grades"
}
