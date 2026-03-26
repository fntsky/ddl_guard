package entity

import "time"

type UserSession struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	UserID int64 `xorm:"not null index 'user_id'"`

	TokenID string `xorm:"not null unique index 'token_id'"`

	RefreshTokenHash string `xorm:"not null 'refresh_token_hash'"`

	ExpiresAt time.Time `xorm:"not null index 'expires_at'"`

	RevokedAt *time.Time `xorm:"null 'revoked_at'"`

	ReplacedByTokenID string `xorm:"null 'replaced_by_token_id'"`

	CreatedAt time.Time `xorm:"created index 'created_at'"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`
}

func (UserSession) TableName() string {
	return "user_session"
}
