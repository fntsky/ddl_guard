package daily_score

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
)

// DailyScoreRepo 平时成绩仓库接口
type DailyScoreRepo interface {
	Create(ctx context.Context, ds *entity.DailyScore) error
	GetByUUID(ctx context.Context, uuid string) (*entity.DailyScore, error)
	GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.DailyScore, error)
	ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.DailyScore, error)
	Update(ctx context.Context, ds *entity.DailyScore) error
	Delete(ctx context.Context, uuid string, userID int64) (int64, error)
	DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error
}

type dailyScoreRepo struct {
	data *data.Data
}

func NewDailyScoreRepo(data *data.Data) DailyScoreRepo {
	return &dailyScoreRepo{data: data}
}

// Create 创建平时成绩
func (r *dailyScoreRepo) Create(ctx context.Context, ds *entity.DailyScore) error {
	_, err := r.data.DB.Context(ctx).Insert(ds)
	return err
}

// GetByUUID 根据UUID获取平时成绩
func (r *dailyScoreRepo) GetByUUID(ctx context.Context, uuid string) (*entity.DailyScore, error) {
	ds := &entity.DailyScore{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(ds)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return ds, nil
}

// GetByUUIDAndUser 根据UUID和用户ID获取平时成绩
func (r *dailyScoreRepo) GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.DailyScore, error) {
	ds := &entity.DailyScore{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(ds)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return ds, nil
}

// ListByFinalGradeID 获取期末成绩下的所有平时成绩
func (r *dailyScoreRepo) ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.DailyScore, error) {
	var dss []*entity.DailyScore
	err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Asc("created_at").
		Find(&dss)
	return dss, err
}

// Update 更新平时成绩
func (r *dailyScoreRepo) Update(ctx context.Context, ds *entity.DailyScore) error {
	_, err := r.data.DB.Context(ctx).
		ID(ds.ID).
		AllCols().
		Update(ds)
	return err
}

// Delete 删除平时成绩
func (r *dailyScoreRepo) Delete(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Delete(&entity.DailyScore{})
}

// DeleteByFinalGradeID 删除期末成绩下的所有平时成绩
func (r *dailyScoreRepo) DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error {
	_, err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Delete(&entity.DailyScore{})
	return err
}
