package final_grade

import (
	"context"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/repo/daily_score"
	"github.com/fntsky/ddl_guard/internal/repo/final_grade"
	"github.com/fntsky/ddl_guard/internal/schema"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

// FinalGradeService 期末成绩服务
type FinalGradeService struct {
	repo         final_grade.FinalGradeRepo
	dailyScoreRepo daily_score.DailyScoreRepo
}

func NewFinalGradeService(repo final_grade.FinalGradeRepo, dailyScoreRepo daily_score.DailyScoreRepo) *FinalGradeService {
	return &FinalGradeService{
		repo:           repo,
		dailyScoreRepo: dailyScoreRepo,
	}
}

// CreateFinalGrade 创建期末成绩
func (s *FinalGradeService) CreateFinalGrade(ctx context.Context, req *schema.CreateFinalGradeReq, userUUID string) (*schema.CreateFinalGradeResp, error) {
	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	// 设置默认比例
	examRatio := req.ExamRatio
	dailyRatio := req.DailyRatio
	if examRatio == 0 && dailyRatio == 0 {
		examRatio = 40
		dailyRatio = 60
	}

	// 验证比例之和为100
	if examRatio+dailyRatio != 100 {
		return nil, apperrors.New(400, apperrors.CodeBadRequest, "exam_ratio + daily_ratio must equal 100")
	}

	now := stime.GetCurrentTime()
	fg := &entity.FinalGrade{
		UUID:       uuid.GenerateUUID(),
		UserID:     userID,
		Name:       strings.TrimSpace(req.Name),
		ExamRatio:  examRatio,
		DailyRatio: dailyRatio,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, fg); err != nil {
		return nil, err
	}

	return &schema.CreateFinalGradeResp{
		UUID:       fg.UUID,
		Name:       fg.Name,
		ExamScore:  fg.ExamScore,
		ExamRatio:  fg.ExamRatio,
		DailyRatio: fg.DailyRatio,
		FinalScore: fg.FinalScore,
	}, nil
}

// GetFinalGrade 获取期末成绩详情
func (s *FinalGradeService) GetFinalGrade(ctx context.Context, fgUUID string, userUUID string) (*schema.FinalGradeDetailResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	fg, err := s.repo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, err
	}
	if fg == nil {
		return nil, apperrors.ErrFinalGradeNotFound
	}

	// 获取关联的平时成绩
	dailyScores, err := s.dailyScoreRepo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, err
	}

	dailyScoreItems := make([]schema.DailyScoreItem, 0, len(dailyScores))
	for _, ds := range dailyScores {
		dailyScoreItems = append(dailyScoreItems, schema.DailyScoreItem{
			UUID:  ds.UUID,
			Type:  string(ds.Type),
			Name:  ds.Name,
			Score: ds.Score,
			Ratio: ds.Ratio,
		})
	}

	return &schema.FinalGradeDetailResp{
		UUID:        fg.UUID,
		Name:        fg.Name,
		ExamScore:   fg.ExamScore,
		ExamRatio:   fg.ExamRatio,
		DailyRatio:  fg.DailyRatio,
		FinalScore:  fg.FinalScore,
		CreatedAt:   fg.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		UpdatedAt:   fg.UpdatedAt.Format("2006-01-02T15:04:05+08:00"),
		DailyScores: dailyScoreItems,
	}, nil
}

// ListFinalGrades 获取期末成绩列表
func (s *FinalGradeService) ListFinalGrades(ctx context.Context, userUUID string, pageReq *schema.PageReq) (*schema.FinalGradeListResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	pageReq.Normalize()

	total, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	fgs, err := s.repo.ListByUserID(ctx, userID, pageReq.Offset(), pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]schema.FinalGradeListItem, 0, len(fgs))
	for _, fg := range fgs {
		list = append(list, schema.FinalGradeListItem{
			UUID:       fg.UUID,
			Name:       fg.Name,
			ExamScore:  fg.ExamScore,
			ExamRatio:  fg.ExamRatio,
			DailyRatio: fg.DailyRatio,
			FinalScore: fg.FinalScore,
			CreatedAt:  fg.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		})
	}

	return &schema.FinalGradeListResp{
		List:     list,
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

// UpdateFinalGrade 更新期末成绩
func (s *FinalGradeService) UpdateFinalGrade(ctx context.Context, fgUUID string, req *schema.UpdateFinalGradeReq, userUUID string) (*schema.UpdateFinalGradeResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	fg, err := s.repo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, err
	}
	if fg == nil {
		return nil, apperrors.ErrFinalGradeNotFound
	}

	// 更新字段
	if req.Name != nil {
		fg.Name = strings.TrimSpace(*req.Name)
	}
	if req.ExamScore != nil {
		fg.ExamScore = *req.ExamScore
	}
	if req.FinalScore != nil {
		fg.FinalScore = *req.FinalScore
	}

	// 处理比例更新
	if req.ExamRatio != nil || req.DailyRatio != nil {
		examRatio := fg.ExamRatio
		dailyRatio := fg.DailyRatio
		if req.ExamRatio != nil {
			examRatio = *req.ExamRatio
		}
		if req.DailyRatio != nil {
			dailyRatio = *req.DailyRatio
		}
		if examRatio+dailyRatio != 100 {
			return nil, apperrors.New(400, apperrors.CodeBadRequest, "exam_ratio + daily_ratio must equal 100")
		}
		fg.ExamRatio = examRatio
		fg.DailyRatio = dailyRatio
	}

	fg.UpdatedAt = stime.GetCurrentTime()

	if err := s.repo.Update(ctx, fg); err != nil {
		return nil, err
	}

	return &schema.UpdateFinalGradeResp{
		UUID:       fg.UUID,
		Name:       fg.Name,
		ExamScore:  fg.ExamScore,
		ExamRatio:  fg.ExamRatio,
		DailyRatio: fg.DailyRatio,
		FinalScore: fg.FinalScore,
	}, nil
}

// DeleteFinalGrade 删除期末成绩
func (s *FinalGradeService) DeleteFinalGrade(ctx context.Context, fgUUID string, userUUID string) error {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return err
	}

	fg, err := s.repo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return err
	}
	if fg == nil {
		return apperrors.ErrFinalGradeNotFound
	}

	// 删除关联的平时成绩
	if err := s.dailyScoreRepo.DeleteByFinalGradeID(ctx, fg.ID); err != nil {
		return err
	}

	affected, err := s.repo.Delete(ctx, fgUUID, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperrors.ErrFinalGradeNotFound
	}

	return nil
}
