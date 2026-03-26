package session

import (
	"context"
	"fmt"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	stime "github.com/fntsky/ddl_guard/pkg/time"
)

type sessionRepo struct {
	data *data.Data
}

func NewSessionRepo(data *data.Data) authsvc.SessionRepo {
	return &sessionRepo{data: data}
}

func (r *sessionRepo) CreateSession(ctx context.Context, session *entity.UserSession) error {
	_, err := r.data.DB.Context(ctx).Insert(session)
	return err
}

func (r *sessionRepo) GetByTokenID(ctx context.Context, tokenID string) (*entity.UserSession, bool, error) {
	session := &entity.UserSession{TokenID: tokenID}
	has, err := r.data.DB.Context(ctx).Get(session)
	return session, has, err
}

func (r *sessionRepo) RotateSession(ctx context.Context, currentTokenID string, newSession *entity.UserSession) error {
	sess := r.data.DB.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	revokedAt := stime.GetCurrentTime()
	affected, err := sess.Context(ctx).
		Where("token_id = ? AND revoked_at IS NULL", currentTokenID).
		Cols("revoked_at", "replaced_by_token_id", "updated_at").
		Update(&entity.UserSession{
			RevokedAt:         &revokedAt,
			ReplacedByTokenID: newSession.TokenID,
			UpdatedAt:         stime.GetCurrentTime(),
		})
	if err != nil {
		_ = sess.Rollback()
		return err
	}
	if affected == 0 {
		_ = sess.Rollback()
		return authsvc.ErrRefreshTokenRevoked
	}

	if _, err = sess.Context(ctx).Insert(newSession); err != nil {
		_ = sess.Rollback()
		return err
	}
	if err = sess.Commit(); err != nil {
		_ = sess.Rollback()
		return fmt.Errorf("commit tx failed: %w", err)
	}
	return nil
}

func (r *sessionRepo) RevokeByTokenID(ctx context.Context, tokenID string) error {
	now := time.Now()
	_, err := r.data.DB.Context(ctx).
		Where("token_id = ? AND revoked_at IS NULL", tokenID).
		Cols("revoked_at", "updated_at").
		Update(&entity.UserSession{
			RevokedAt: &now,
			UpdatedAt: stime.GetCurrentTime(),
		})
	return err
}
