package quiz_score

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
)

// QuizScoreRepo 小测成绩仓库接口
type QuizScoreRepo interface {
	Create(ctx context.Context, qs *entity.QuizScore) error
	GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.QuizScore, error)
	ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.QuizScore, error)
	GetAvgByFinalGradeID(ctx context.Context, finalGradeID int64) (float64, error)
	Update(ctx context.Context, qs *entity.QuizScore) error
	Delete(ctx context.Context, uuid string, userID int64) (int64, error)
	DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error
}

type quizScoreRepo struct {
	data *data.Data
}

func NewQuizScoreRepo(data *data.Data) QuizScoreRepo {
	return &quizScoreRepo{data: data}
}

func (r *quizScoreRepo) Create(ctx context.Context, qs *entity.QuizScore) error {
	_, err := r.data.DB.Context(ctx).Insert(qs)
	return err
}

func (r *quizScoreRepo) GetByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.QuizScore, error) {
	qs := &entity.QuizScore{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(qs)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return qs, nil
}

func (r *quizScoreRepo) ListByFinalGradeID(ctx context.Context, finalGradeID int64) ([]*entity.QuizScore, error) {
	var list []*entity.QuizScore
	err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Asc("created_at").
		Find(&list)
	return list, err
}

func (r *quizScoreRepo) GetAvgByFinalGradeID(ctx context.Context, finalGradeID int64) (float64, error) {
	var list []*entity.QuizScore
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
	for _, qs := range list {
		sum += qs.Score
	}
	return sum / float64(len(list)), nil
}

func (r *quizScoreRepo) Update(ctx context.Context, qs *entity.QuizScore) error {
	_, err := r.data.DB.Context(ctx).
		ID(qs.ID).
		AllCols().
		Update(qs)
	return err
}

func (r *quizScoreRepo) Delete(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Delete(&entity.QuizScore{})
}

func (r *quizScoreRepo) DeleteByFinalGradeID(ctx context.Context, finalGradeID int64) error {
	_, err := r.data.DB.Context(ctx).
		Where("final_grade_id = ?", finalGradeID).
		Delete(&entity.QuizScore{})
	return err
}