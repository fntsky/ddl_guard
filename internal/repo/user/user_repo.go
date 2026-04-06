package user

import (
	"context"

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

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.data.DB.Context(ctx).
		Where("email = ?", email).
		Exist(&entity.User{})
}

func (r *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := r.data.DB.Context(ctx).Insert(user)
	return err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	has, err := r.data.DB.Context(ctx).
		Where("email = ?", email).
		Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}
