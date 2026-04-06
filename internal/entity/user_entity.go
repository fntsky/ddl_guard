package entity

import "time"

type User struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	UUID string `xorm:"not null unique index 'uuid'"`

	Username string `xorm:"not null index 'username'"`

	Email string `xorm:"null unique 'email'"`

	PasswordHash string `xorm:"not null 'password_hash'"`

	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`
}

func (User) TableName() string {
	return "user"
}
