package user_auth

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
	usersvc "github.com/fntsky/ddl_guard/internal/service/user"
)

type userAuthRepo struct {
	data *data.Data
}

func NewUserAuthRepo(data *data.Data) usersvc.UserAuthRepo {
	return &userAuthRepo{data: data}
}

func (r *userAuthRepo) GetByTypeAndIdentifier(ctx context.Context, authType, identifier string) (*entity.UserAuth, error) {
	auth := &entity.UserAuth{}
	has, err := r.data.DB.Context(ctx).
		Where("auth_type = ? AND auth_identifier = ?", authType, identifier).
		Get(auth)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return auth, nil
}

func (r *userAuthRepo) Create(ctx context.Context, auth *entity.UserAuth) error {
	_, err := r.data.DB.Context(ctx).Insert(auth)
	return err
}
