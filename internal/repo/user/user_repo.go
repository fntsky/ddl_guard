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

func (r *userRepo) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	_, err := r.data.DB.Context(ctx).
		ID(userID).
		Update(&entity.User{PasswordHash: passwordHash})
	return err
}

func (r *userRepo) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	user := &entity.User{}
	has, err := r.data.DB.Context(ctx).ID(userID).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

// GetUserEmailsByIDs 批量获取用户邮箱，返回 userID -> email 映射
func (r *userRepo) GetUserEmailsByIDs(ctx context.Context, userIDs []int64) (map[int64]string, error) {
	result := make(map[int64]string)
	if len(userIDs) == 0 {
		return result, nil
	}

	var users []struct {
		ID    int64  `xorm:"'id'"`
		Email string `xorm:"'email'"`
	}
	err := r.data.DB.Context(ctx).
		Table("user").
		In("id", userIDs).
		Where("email IS NOT NULL").
		And("email != ''").
		Cols("id", "email").
		Find(&users)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		result[u.ID] = u.Email
	}
	return result, nil
}
