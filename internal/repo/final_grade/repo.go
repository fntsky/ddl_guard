package final_grade

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
)

// FinalGradeRepo 期末成绩仓库接口
type FinalGradeRepo interface {
	Create(ctx context.Context, fg *entity.FinalGrade) error
	GetByUUID(ctx context.Context, uuid string) (*entity.FinalGrade, error)
	GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.FinalGrade, error)
	GetByID(ctx context.Context, id int64) (*entity.FinalGrade, error)
	ListByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.FinalGrade, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	Update(ctx context.Context, fg *entity.FinalGrade) error
	Delete(ctx context.Context, uuid string, userID int64) (int64, error)
	GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error)
}

type finalGradeRepo struct {
	data *data.Data
}

func NewFinalGradeRepo(data *data.Data) FinalGradeRepo {
	return &finalGradeRepo{data: data}
}

// Create 创建期末成绩
func (r *finalGradeRepo) Create(ctx context.Context, fg *entity.FinalGrade) error {
	_, err := r.data.DB.Context(ctx).Insert(fg)
	return err
}

// GetByUUID 根据UUID获取期末成绩
func (r *finalGradeRepo) GetByUUID(ctx context.Context, uuid string) (*entity.FinalGrade, error) {
	fg := &entity.FinalGrade{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(fg)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return fg, nil
}

// GetByUUIDAndUser 根据UUID和用户ID获取期末成绩
func (r *finalGradeRepo) GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.FinalGrade, error) {
	fg := &entity.FinalGrade{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(fg)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return fg, nil
}

// ListByUserID 分页获取用户的期末成绩列表
func (r *finalGradeRepo) GetByID(ctx context.Context, id int64) (*entity.FinalGrade, error) {
	fg := &entity.FinalGrade{}
	has, err := r.data.DB.Context(ctx).ID(id).Get(fg)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return fg, nil
}

func (r *finalGradeRepo) ListByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.FinalGrade, error) {
	var fgs []*entity.FinalGrade
	err := r.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		Desc("created_at").
		Limit(limit, offset).
		Find(&fgs)
	return fgs, err
}

// CountByUserID 统计用户的期末成绩数量
func (r *finalGradeRepo) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		Count(&entity.FinalGrade{})
}

// Update 更新期末成绩
func (r *finalGradeRepo) Update(ctx context.Context, fg *entity.FinalGrade) error {
	_, err := r.data.DB.Context(ctx).
		ID(fg.ID).
		AllCols().
		Update(fg)
	return err
}

// Delete 删除期末成绩
func (r *finalGradeRepo) Delete(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Delete(&entity.FinalGrade{})
}

// GetUserIDByUserUUID 根据用户UUID获取用户ID
func (r *finalGradeRepo) GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error) {
	user := &entity.User{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, apperrors.ErrUserNotFound
	}
	return user.ID, nil
}
