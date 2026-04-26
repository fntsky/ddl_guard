package exam

import (
	"context"
	"errors"

	"github.com/fntsky/ddl_guard/internal/base/data"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
)

// ExamRepo 考试仓库接口（在service层定义）
type ExamRepo interface {
	Create(ctx context.Context, exam *entity.Exam) error
	GetByUUID(ctx context.Context, uuid string) (*entity.Exam, error)
	GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.Exam, error)
	ListByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.Exam, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	Update(ctx context.Context, exam *entity.Exam) error
	Delete(ctx context.Context, uuid string, userID int64) (int64, error)
	GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error)
}

type examRepo struct {
	data *data.Data
}

func NewExamRepo(data *data.Data) ExamRepo {
	return &examRepo{data: data}
}

// Create 创建考试
func (r *examRepo) Create(ctx context.Context, exam *entity.Exam) error {
	_, err := r.data.DB.Context(ctx).Insert(exam)
	return err
}

// GetByUUID 根据UUID获取考试
func (r *examRepo) GetByUUID(ctx context.Context, uuid string) (*entity.Exam, error) {
	exam := &entity.Exam{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(exam)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return exam, nil
}

// GetByUUIDAndUser 根据UUID和用户ID获取考试
func (r *examRepo) GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.Exam, error) {
	exam := &entity.Exam{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(exam)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return exam, nil
}

// ListByUserID 分页获取用户的考试列表
func (r *examRepo) ListByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.Exam, error) {
	var exams []*entity.Exam
	err := r.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		Desc("start_time").
		Limit(limit, offset).
		Find(&exams)
	return exams, err
}

// CountByUserID 统计用户的考试数量
func (r *examRepo) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		Count(&entity.Exam{})
}

// Update 更新考试
func (r *examRepo) Update(ctx context.Context, exam *entity.Exam) error {
	_, err := r.data.DB.Context(ctx).
		ID(exam.ID).
		AllCols().
		Update(exam)
	return err
}

// Delete 删除考试
func (r *examRepo) Delete(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Delete(&entity.Exam{})
}

// GetUserIDByUserUUID 根据用户UUID获取用户ID
func (r *examRepo) GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error) {
	user := &entity.User{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, apperrors.ErrUserNotFound
	}
	if user.ID <= 0 {
		return 0, errors.New("invalid user id")
	}
	return user.ID, nil
}
