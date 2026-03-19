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
	UUID            string    `xorm:"varchar(36) not null pk 'uuid'"`
	CreatedAt       time.Time `xorm:"created 'created_at'"`
	UpdatedAt       time.Time `xorm:"updated 'updated_at'"`
	DeadLine        time.Time `xorm:"'deadline' not null"`
	Status          int       `xorm:"not null default 0 index 'status'"`
	EealyRemindTime time.Time `xorm:"'early_remind_time'"`
	Title           string    `xorm:"varchar(255) not null 'title'"`
	Description     string    `xorm:"text 'description'"`
}

func (DDL) TableName() string {
	return "ddl"
}
