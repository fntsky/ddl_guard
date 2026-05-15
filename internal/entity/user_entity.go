package entity

import "time"

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type User struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	UUID string `xorm:"not null unique index 'uuid'"`

	Username string `xorm:"not null index 'username'"`

	Email *string `xorm:"null unique 'email'"`

	Phone *string `xorm:"null unique 'phone'"`

	PasswordHash string `xorm:"not null 'password_hash'"`

	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`
}

func (User) TableName() string {
	return "user"
}
