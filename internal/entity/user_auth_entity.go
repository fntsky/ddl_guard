package entity

import "time"

const (
	UserAuthTypeEmail = "email"
)

type UserAuth struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	UserID int64 `xorm:"not null index 'user_id'"`

	AuthType string `xorm:"not null unique(idx_auth_type_identifier)"`

	AuthIdentifier string `xorm:"not null unique(idx_auth_type_identifier)"`

	AuthMeta string `xorm:"null 'auth_meta'"`

	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`
}

func (UserAuth) TableName() string {
	return "user_auth"
}
