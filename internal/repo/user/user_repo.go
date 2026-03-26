package user

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
	usersvc "github.com/fntsky/ddl_guard/internal/service/user"
)

type userRepo struct {
	data *data.Data
}

func NewUserRepo(data *data.Data) usersvc.UserRepo {
	return &userRepo{data: data}
}

func (r *userRepo) ExistsByAuthIdentifier(ctx context.Context, authType string, authIdentifier string) (bool, error) {
	item := &entity.UserAuth{}
	return r.data.DB.Context(ctx).
		Where("auth_type = ? AND auth_identifier = ?", authType, authIdentifier).
		Get(item)
}

func (r *userRepo) CreateUserWithAuth(ctx context.Context, user *entity.User, auth *entity.UserAuth) error {
	session := r.data.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	if _, err := session.Context(ctx).Insert(user); err != nil {
		_ = session.Rollback()
		return err
	}
	auth.UserID = user.ID
	if _, err := session.Context(ctx).Insert(auth); err != nil {
		_ = session.Rollback()
		return err
	}
	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		return fmt.Errorf("commit tx failed: %w", err)
	}
	return nil
}
