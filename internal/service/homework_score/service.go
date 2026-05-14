package homework_score

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/repo/final_grade"
	"github.com/fntsky/ddl_guard/internal/repo/homework_score"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/pkg/uuid"

	fg_service "github.com/fntsky/ddl_guard/internal/service/final_grade"
)

type HomeworkScoreService struct {
	homeworkScoreRepo homework_score.HomeworkScoreRepo
	finalGradeRepo    final_grade.FinalGradeRepo
	finalGradeSvc     *fg_service.FinalGradeService
}

func NewHomeworkScoreService(
	homeworkScoreRepo homework_score.HomeworkScoreRepo,
	finalGradeRepo final_grade.FinalGradeRepo,
	finalGradeSvc *fg_service.FinalGradeService,
) *HomeworkScoreService {
	return &HomeworkScoreService{
		homeworkScoreRepo: homeworkScoreRepo,
		finalGradeRepo:    finalGradeRepo,
		finalGradeSvc:     finalGradeSvc,
	}
}

func (s *HomeworkScoreService) CreateHomeworkScore(ctx context.Context, fgUUID string, req *schema.CreateHomeworkScoreReq, userUUID string) (*schema.CreateHomeworkScoreResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get final grade failed: %w", err)
	}
	if fg == nil {
		return nil, fmt.Errorf("final grade not found")
	}

	hs := &entity.HomeworkScore{
		UUID:         uuid.GenerateUUID(),
		FinalGradeID: fg.ID,
		UserID:       userID,
		Name:         req.Name,
		Score:        req.Score,
	}

	if err := s.homeworkScoreRepo.Create(ctx, hs); err != nil {
		return nil, fmt.Errorf("create homework score failed: %w", err)
	}

	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return nil, fmt.Errorf("recalculate final score failed: %w", err)
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil, fmt.Errorf("update final grade failed: %w", err)
	}

	return &schema.CreateHomeworkScoreResp{
		UUID:         hs.UUID,
		FinalGradeID: fgUUID,
		Name:         hs.Name,
		Score:        hs.Score,
	}, nil
}

func (s *HomeworkScoreService) ListHomeworkScores(ctx context.Context, fgUUID string, userUUID string, pageReq *schema.PageReq) (*schema.HomeworkScoreListResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get final grade failed: %w", err)
	}
	if fg == nil {
		return nil, fmt.Errorf("final grade not found")
	}

	list, err := s.homeworkScoreRepo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, fmt.Errorf("list homework scores failed: %w", err)
	}

	items := make([]schema.HomeworkScoreListItem, 0, len(list))
	for _, hs := range list {
		items = append(items, schema.HomeworkScoreListItem{
			UUID:     hs.UUID,
			Name:     hs.Name,
			Score:    hs.Score,
			CreateAt: hs.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		})
	}

	return &schema.HomeworkScoreListResp{
		List:     items,
		Total:    int64(len(items)),
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

func (s *HomeworkScoreService) UpdateHomeworkScore(ctx context.Context, hsUUID string, req *schema.UpdateHomeworkScoreReq, userUUID string) (*schema.UpdateHomeworkScoreResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	hs, err := s.homeworkScoreRepo.GetByUUIDAndUser(ctx, hsUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get homework score failed: %w", err)
	}
	if hs == nil {
		return nil, fmt.Errorf("homework score not found")
	}

	if req.Name != nil {
		hs.Name = *req.Name
	}
	if req.Score != nil {
		hs.Score = *req.Score
	}

	if err := s.homeworkScoreRepo.Update(ctx, hs); err != nil {
		return nil, fmt.Errorf("update homework score failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByID(ctx, hs.FinalGradeID)
	if err != nil {
		return nil, fmt.Errorf("get final grade failed: %w", err)
	}
	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return nil, fmt.Errorf("recalculate final score failed: %w", err)
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil, fmt.Errorf("update final grade failed: %w", err)
	}

	return &schema.UpdateHomeworkScoreResp{
		UUID:  hs.UUID,
		Name:  hs.Name,
		Score: hs.Score,
	}, nil
}

func (s *HomeworkScoreService) DeleteHomeworkScore(ctx context.Context, hsUUID string, userUUID string) error {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("get user id failed: %w", err)
	}

	hs, err := s.homeworkScoreRepo.GetByUUIDAndUser(ctx, hsUUID, userID)
	if err != nil {
		return fmt.Errorf("get homework score failed: %w", err)
	}
	if hs == nil {
		return fmt.Errorf("homework score not found")
	}

	if _, err := s.homeworkScoreRepo.Delete(ctx, hsUUID, userID); err != nil {
		return fmt.Errorf("delete homework score failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByID(ctx, hs.FinalGradeID)
	if err != nil {
		return fmt.Errorf("get final grade failed: %w", err)
	}
	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return nil
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil
	}

	return nil
}
