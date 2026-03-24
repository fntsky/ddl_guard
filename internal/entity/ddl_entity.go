package entity

import "time"

const (
	DDLStatusDraft   = 0
	DDLStatusActive  = 1
	DDLStatusExpired = 2
	DDLStatusDeleted = 3
	DDLStatusDone    = 4
)

type DDL struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	
	UUID string `xorm:"uuid not null unique index 'uuid'"`

	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`

	DeadLine time.Time `xorm:"not null index 'idx_status_deadline' 'deadline'"`

	Status int `xorm:"not null default 0 index 'idx_status_deadline' 'status'"`

	EealyRemindTime time.Time `xorm:"index 'early_remind_time'"`

	Title string `xorm:"text not null 'title'"`

	Description string `xorm:"text 'description'"`
}

func (DDL) TableName() string {
	return "ddl"
}
