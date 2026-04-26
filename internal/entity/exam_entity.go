package entity

import "time"

type Exam struct {
	ID        int64     `xorm:"pk autoincr 'id'"`
	UUID      string    `xorm:"uuid not null unique index 'uuid'"`
	UserID    int64     `xorm:"not null index 'user_id'"`
	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`
	Name      string    `xorm:"text not null 'name'"`
	StartTime time.Time `xorm:"not null index 'start_time'"`
	EndTime   time.Time `xorm:"not null index 'end_time'"`
	Location  string    `xorm:"text 'location'"`
	Notes     string    `xorm:"text 'notes'"`
}

func (Exam) TableName() string {
	return "exam"
}
