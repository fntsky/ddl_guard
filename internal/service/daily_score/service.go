package daily_score

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

// DailyScoreService 平时成绩服务
type DailyScoreService struct {
	repo           daily_score.DailyScoreRepo
	finalGradeRepo final_grade.FinalGradeRepo
}

func NewDailyScoreService(repo daily_score.DailyScoreRepo, finalGradeRepo final_grade.FinalGradeRepo) *DailyScoreService {
	return &DailyScoreService{
		repo:           repo,
		finalGradeRepo: finalGradeRepo,
	}
}

// CreateDailyScore 创建平时成绩
func (s *DailyScoreService) CreateDailyScore(ctx context.Context, fgUUID string, req *schema.CreateDailyScoreReq, userUUID string) (*schema.CreateDailyScoreResp, error) {
	// 获取用户ID
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	// 获取期末成绩记录
	fg, err := s.finalGradeRepo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, err
	}
	if fg == nil {
		return nil, apperrors.ErrFinalGradeNotFound
	}

	// 验证分数范围
	if req.Score < 0 || req.Score > 100 {
		return nil, apperrors.New(400, apperrors.CodeBadRequest, "score must be between 0 and 100")
	}

	// 验证比例范围
	if req.Ratio < 0 || req.Ratio > 100 {
		return nil, apperrors.New(400, apperrors.CodeBadRequest, "ratio must be between 0 and 100")
	}

	now := stime.GetCurrentTime()
	ds := &entity.DailyScore{
		UUID:         uuid.GenerateUUID(),
		FinalGradeID: fg.ID,
		UserID:       userID,
		Type:         entity.DailyScoreType(req.Type),
		Name:         strings.TrimSpace(req.Name),
		Score:        req.Score,
		Ratio:        req.Ratio,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, ds); err != nil {
		return nil, err
	}

	return &schema.CreateDailyScoreResp{
		UUID:         ds.UUID,
		FinalGradeID: fg.UUID,
		Type:         string(ds.Type),
		Name:         ds.Name,
		Score:        ds.Score,
		Ratio:        ds.Ratio,
	}, nil
}

// ListDailyScores 获取平时成绩列表
func (s *DailyScoreService) ListDailyScores(ctx context.Context, fgUUID string, userUUID string) (*schema.DailyScoreListResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	fg, err := s.finalGradeRepo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, err
	}
	if fg == nil {
		return nil, apperrors.ErrFinalGradeNotFound
	}

	dailyScores, err := s.repo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, err
	}

	list := make([]schema.DailyScoreListItem, 0, len(dailyScores))
	for _, ds := range dailyScores {
		list = append(list, schema.DailyScoreListItem{
			UUID:     ds.UUID,
			Type:     string(ds.Type),
			Name:     ds.Name,
			Score:    ds.Score,
			Ratio:    ds.Ratio,
			CreateAt: ds.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		})
	}

	return &schema.DailyScoreListResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

// UpdateDailyScore 更新平时成绩
func (s *DailyScoreService) UpdateDailyScore(ctx context.Context, dsUUID string, req *schema.UpdateDailyScoreReq, userUUID string) (*schema.UpdateDailyScoreResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	ds, err := s.repo.GetByUUIDAndUser(ctx, dsUUID, userID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, apperrors.ErrDailyScoreNotFound
	}

	// 更新字段
	if req.Type != nil {
		ds.Type = entity.DailyScoreType(*req.Type)
	}
	if req.Name != nil {
		ds.Name = strings.TrimSpace(*req.Name)
	}
	if req.Score != nil {
		if *req.Score < 0 || *req.Score > 100 {
			return nil, apperrors.New(400, apperrors.CodeBadRequest, "score must be between 0 and 100")
		}
		ds.Score = *req.Score
	}
	if req.Ratio != nil {
		if *req.Ratio < 0 || *req.Ratio > 100 {
			return nil, apperrors.New(400, apperrors.CodeBadRequest, "ratio must be between 0 and 100")
		}
		ds.Ratio = *req.Ratio
	}

	ds.UpdatedAt = stime.GetCurrentTime()

	if err := s.repo.Update(ctx, ds); err != nil {
		return nil, err
	}

	return &schema.UpdateDailyScoreResp{
		UUID:  ds.UUID,
		Type:  string(ds.Type),
		Name:  ds.Name,
		Score: ds.Score,
		Ratio: ds.Ratio,
	}, nil
}

// DeleteDailyScore 删除平时成绩
func (s *DailyScoreService) DeleteDailyScore(ctx context.Context, dsUUID string, userUUID string) error {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return err
	}

	affected, err := s.repo.Delete(ctx, dsUUID, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperrors.ErrDailyScoreNotFound
	}

	return nil
}
