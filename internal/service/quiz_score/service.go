package quiz_score

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/repo/final_grade"
	"github.com/fntsky/ddl_guard/internal/repo/quiz_score"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/pkg/uuid"

	fg_service "github.com/fntsky/ddl_guard/internal/service/final_grade"
)

type QuizScoreService struct {
	quizScoreRepo   quiz_score.QuizScoreRepo
	finalGradeRepo  final_grade.FinalGradeRepo
	finalGradeSvc   *fg_service.FinalGradeService
}

func NewQuizScoreService(
	quizScoreRepo quiz_score.QuizScoreRepo,
	finalGradeRepo final_grade.FinalGradeRepo,
	finalGradeSvc *fg_service.FinalGradeService,
) *QuizScoreService {
	return &QuizScoreService{
		quizScoreRepo:  quizScoreRepo,
		finalGradeRepo: finalGradeRepo,
		finalGradeSvc:  finalGradeSvc,
	}
}

func (s *QuizScoreService) CreateQuizScore(ctx context.Context, fgUUID string, req *schema.CreateQuizScoreReq, userUUID string) (*schema.CreateQuizScoreResp, error) {
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

	qs := &entity.QuizScore{
		UUID:         uuid.GenerateUUID(),
		FinalGradeID: fg.ID,
		UserID:       userID,
		Name:         req.Name,
		Score:        req.Score,
	}

	if err := s.quizScoreRepo.Create(ctx, qs); err != nil {
		return nil, fmt.Errorf("create quiz score failed: %w", err)
	}

	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return nil, fmt.Errorf("recalculate final score failed: %w", err)
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil, fmt.Errorf("update final grade failed: %w", err)
	}

	return &schema.CreateQuizScoreResp{
		UUID:         qs.UUID,
		FinalGradeID: fgUUID,
		Name:         qs.Name,
		Score:        qs.Score,
	}, nil
}

func (s *QuizScoreService) ListQuizScores(ctx context.Context, fgUUID string, userUUID string, pageReq *schema.PageReq) (*schema.QuizScoreListResp, error) {
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

	list, err := s.quizScoreRepo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, fmt.Errorf("list quiz scores failed: %w", err)
	}

	items := make([]schema.QuizScoreListItem, 0, len(list))
	for _, qs := range list {
		items = append(items, schema.QuizScoreListItem{
			UUID:     qs.UUID,
			Name:     qs.Name,
			Score:    qs.Score,
			CreateAt: qs.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		})
	}

	return &schema.QuizScoreListResp{
		List:     items,
		Total:    int64(len(items)),
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

func (s *QuizScoreService) UpdateQuizScore(ctx context.Context, qsUUID string, req *schema.UpdateQuizScoreReq, userUUID string) (*schema.UpdateQuizScoreResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	qs, err := s.quizScoreRepo.GetByUUIDAndUser(ctx, qsUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get quiz score failed: %w", err)
	}
	if qs == nil {
		return nil, fmt.Errorf("quiz score not found")
	}

	if req.Name != nil {
		qs.Name = *req.Name
	}
	if req.Score != nil {
		qs.Score = *req.Score
	}

	if err := s.quizScoreRepo.Update(ctx, qs); err != nil {
		return nil, fmt.Errorf("update quiz score failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByID(ctx, qs.FinalGradeID)
	if err != nil {
		return nil, fmt.Errorf("get final grade failed: %w", err)
	}
	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return nil, fmt.Errorf("recalculate final score failed: %w", err)
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil, fmt.Errorf("update final grade failed: %w", err)
	}

	return &schema.UpdateQuizScoreResp{
		UUID:  qs.UUID,
		Name:  qs.Name,
		Score: qs.Score,
	}, nil
}

func (s *QuizScoreService) DeleteQuizScore(ctx context.Context, qsUUID string, userUUID string) error {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("get user id failed: %w", err)
	}

	qs, err := s.quizScoreRepo.GetByUUIDAndUser(ctx, qsUUID, userID)
	if err != nil {
		return fmt.Errorf("get quiz score failed: %w", err)
	}
	if qs == nil {
		return fmt.Errorf("quiz score not found")
	}

	if _, err := s.quizScoreRepo.Delete(ctx, qsUUID, userID); err != nil {
		return fmt.Errorf("delete quiz score failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByID(ctx, qs.FinalGradeID)
	if err != nil {
		return fmt.Errorf("get final grade failed: %w", err)
	}
	if err := s.finalGradeSvc.RecalculateFinalScore(ctx, fg); err != nil {
		return fmt.Errorf("recalculate final score failed: %w", err)
	}
	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil
	}

	return nil
}
