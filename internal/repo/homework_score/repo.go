package homework_score

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
)

// HomeworkScoreRepo 作业成绩仓库接口
type HomeworkScoreRepo interface {
	Create(ctx context.Context, hs *entity.HomeworkScore) error
	GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.HomeworkScore, error)
	ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.HomeworkScore, error)
	GetAvgByFinalGradeID(ctx context.Context, finalGradeID int64) (float64, error)
	Update(ctx context.Context, hs *entity.HomeworkScore) error
	Delete(ctx context.Context, uuid string, userID int64) (int64, error)
	DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error
}

type homeworkScoreRepo struct {
	data *data.Data
}

func NewHomeworkScoreRepo(data *data.Data) HomeworkScoreRepo {
	return &homeworkScoreRepo{data: data}
}

func (r *homeworkScoreRepo) Create(ctx context.Context, hs *entity.HomeworkScore) error {
	_, err := r.data.DB.Context(ctx).Insert(hs)
	return err
}

func (r *homeworkScoreRepo) GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.HomeworkScore, error) {
	hs := &entity.HomeworkScore{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(hs)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return hs, nil
}

func (r *homeworkScoreRepo) ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.HomeworkScore, error) {
	var list []*entity.HomeworkScore
	err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Asc("created_at").
		Find(&list)
	return list, err
}

func (r *homeworkScoreRepo) GetAvgByFinalGradeID(ctx context.Context, finalGradeID int64) (float64, error) {
	var list []*entity.HomeworkScore
	err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Find(&list)
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, nil
	}
	var sum float64
	for _, hs := range list {
		sum += hs.Score
	}
	return sum / float64(len(list)), nil
}

func (r *homeworkScoreRepo) Update(ctx context.Context, hs *entity.HomeworkScore) error {
	_, err := r.data.DB.Context(ctx).
		ID(hs.ID).
		AllCols().
		Update(hs)
	return err
}

func (r *homeworkScoreRepo) Delete(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Delete(&entity.HomeworkScore{})
}

func (r *homeworkScoreRepo) DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error {
	_, err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Delete(&entity.HomeworkScore{})
	return err
}