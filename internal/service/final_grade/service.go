package final_grade

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/repo/final_grade"
	"github.com/fntsky/ddl_guard/internal/repo/homework_score"
	"github.com/fntsky/ddl_guard/internal/repo/quiz_score"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

type FinalGradeService struct {
	finalGradeRepo   final_grade.FinalGradeRepo
	quizScoreRepo    quiz_score.QuizScoreRepo
	homeworkScoreRepo homework_score.HomeworkScoreRepo
}

func NewFinalGradeService(
	finalGradeRepo final_grade.FinalGradeRepo,
	quizScoreRepo quiz_score.QuizScoreRepo,
	homeworkScoreRepo homework_score.HomeworkScoreRepo,
) *FinalGradeService {
	return &FinalGradeService{
		finalGradeRepo:    finalGradeRepo,
		quizScoreRepo:     quizScoreRepo,
		homeworkScoreRepo: homeworkScoreRepo,
	}
}

func (s *FinalGradeService) CreateFinalGrade(ctx context.Context, req *schema.CreateFinalGradeReq, userUUID string) (*schema.CreateFinalGradeResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	if err := s.validateRatios(req.ExamRatio, req.ClassroomBonusRatio, req.AttendanceRatio, req.QuizRatio, req.HomeworkRatio); err != nil {
		return nil, err
	}

	fg := &entity.FinalGrade{
		UUID:                  uuid.GenerateUUID(),
		UserID:                userID,
		Name:                  req.Name,
		ExamRatio:             req.ExamRatio,
		ClassroomBonusRatio:   req.ClassroomBonusRatio,
		AttendanceRatio:       req.AttendanceRatio,
		QuizRatio:             req.QuizRatio,
		HomeworkRatio:         req.HomeworkRatio,
	}

	if err := s.finalGradeRepo.Create(ctx, fg); err != nil {
		return nil, fmt.Errorf("create final grade failed: %w", err)
	}

	return &schema.CreateFinalGradeResp{
		UUID:                fg.UUID,
		Name:                fg.Name,
		ExamScore:           fg.ExamScore,
		ExamRatio:           fg.ExamRatio,
		ClassroomBonusScore: fg.ClassroomBonusScore,
		ClassroomBonusRatio: fg.ClassroomBonusRatio,
		AttendanceScore:     fg.AttendanceScore,
		AttendanceRatio:     fg.AttendanceRatio,
		QuizRatio:           fg.QuizRatio,
		HomeworkRatio:       fg.HomeworkRatio,
		FinalScore:          fg.FinalScore,
	}, nil
}

func (s *FinalGradeService) ListFinalGrades(ctx context.Context, userUUID string, pageReq *schema.PageReq) (*schema.FinalGradeListResp, error) {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %w", err)
	}

	pageReq.Normalize()

	total, err := s.finalGradeRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("count final grades failed: %w", err)
	}

	list, err := s.finalGradeRepo.ListByUserID(ctx, userID, pageReq.Offset(), pageReq.PageSize)

	items := make([]schema.FinalGradeListItem, 0, len(list))
	for _, fg := range list {
		items = append(items, schema.FinalGradeListItem{
			UUID:                fg.UUID,
			Name:                fg.Name,
			ExamScore:           fg.ExamScore,
			ExamRatio:           fg.ExamRatio,
			ClassroomBonusScore: fg.ClassroomBonusScore,
			ClassroomBonusRatio: fg.ClassroomBonusRatio,
			AttendanceScore:     fg.AttendanceScore,
			AttendanceRatio:     fg.AttendanceRatio,
			QuizRatio:           fg.QuizRatio,
			HomeworkRatio:       fg.HomeworkRatio,
			FinalScore:          fg.FinalScore,
			CreatedAt:           fg.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		})
	}

	return &schema.FinalGradeListResp{
		List:     items,
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

func (s *FinalGradeService) GetFinalGrade(ctx context.Context, fgUUID string, userUUID string) (*schema.FinalGradeDetailResp, error) {
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

	// Get quiz scores
	quizList, err := s.quizScoreRepo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, fmt.Errorf("get quiz scores failed: %w", err)
	}
	quizItems := make([]schema.QuizScoreItem, 0, len(quizList))
	for _, qs := range quizList {
		quizItems = append(quizItems, schema.QuizScoreItem{
			UUID:  qs.UUID,
			Name:  qs.Name,
			Score: qs.Score,
		})
	}

	// Get homework scores
	homeworkList, err := s.homeworkScoreRepo.ListByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return nil, fmt.Errorf("get homework scores failed: %w", err)
	}
	homeworkItems := make([]schema.HomeworkScoreItem, 0, len(homeworkList))
	for _, hs := range homeworkList {
		homeworkItems = append(homeworkItems, schema.HomeworkScoreItem{
			UUID:  hs.UUID,
			Name:  hs.Name,
			Score: hs.Score,
		})
	}

	return &schema.FinalGradeDetailResp{
		UUID:                fg.UUID,
		Name:                fg.Name,
		ExamScore:           fg.ExamScore,
		ExamRatio:           fg.ExamRatio,
		ClassroomBonusScore: fg.ClassroomBonusScore,
		ClassroomBonusRatio: fg.ClassroomBonusRatio,
		AttendanceScore:     fg.AttendanceScore,
		AttendanceRatio:     fg.AttendanceRatio,
		QuizRatio:           fg.QuizRatio,
		HomeworkRatio:       fg.HomeworkRatio,
		FinalScore:          fg.FinalScore,
		CreatedAt:           fg.CreatedAt.Format("2006-01-02T15:04:05+08:00"),
		UpdatedAt:           fg.UpdatedAt.Format("2006-01-02T15:04:05+08:00"),
		QuizScores:          quizItems,
		HomeworkScores:      homeworkItems,
	}, nil
}

func (s *FinalGradeService) UpdateFinalGrade(ctx context.Context, fgUUID string, req *schema.UpdateFinalGradeReq, userUUID string) (*schema.UpdateFinalGradeResp, error) {
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

	if req.Name != nil {
		fg.Name = *req.Name
	}
	if req.ExamScore != nil {
		fg.ExamScore = *req.ExamScore
	}
	if req.ExamRatio != nil {
		fg.ExamRatio = *req.ExamRatio
	}
	if req.ClassroomBonusScore != nil {
		fg.ClassroomBonusScore = *req.ClassroomBonusScore
	}
	if req.ClassroomBonusRatio != nil {
		fg.ClassroomBonusRatio = *req.ClassroomBonusRatio
	}
	if req.AttendanceScore != nil {
		fg.AttendanceScore = *req.AttendanceScore
	}
	if req.AttendanceRatio != nil {
		fg.AttendanceRatio = *req.AttendanceRatio
	}
	if req.QuizRatio != nil {
		fg.QuizRatio = *req.QuizRatio
	}
	if req.HomeworkRatio != nil {
		fg.HomeworkRatio = *req.HomeworkRatio
	}

	if err := s.validateRatios(fg.ExamRatio, fg.ClassroomBonusRatio, fg.AttendanceRatio, fg.QuizRatio, fg.HomeworkRatio); err != nil {
		return nil, err
	}

	if err := s.RecalculateFinalScore(ctx, fg); err != nil {
		return nil, fmt.Errorf("recalculate final score failed: %w", err)
	}

	if err := s.finalGradeRepo.Update(ctx, fg); err != nil {
		return nil, fmt.Errorf("update final grade failed: %w", err)
	}

	return &schema.UpdateFinalGradeResp{
		UUID:                fg.UUID,
		Name:                fg.Name,
		ExamScore:           fg.ExamScore,
		ExamRatio:           fg.ExamRatio,
		ClassroomBonusScore: fg.ClassroomBonusScore,
		ClassroomBonusRatio: fg.ClassroomBonusRatio,
		AttendanceScore:     fg.AttendanceScore,
		AttendanceRatio:     fg.AttendanceRatio,
		QuizRatio:           fg.QuizRatio,
		HomeworkRatio:       fg.HomeworkRatio,
		FinalScore:          fg.FinalScore,
	}, nil
}

func (s *FinalGradeService) DeleteFinalGrade(ctx context.Context, fgUUID string, userUUID string) error {
	userID, err := s.finalGradeRepo.GetUserIDByUserUUID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("get user id failed: %w", err)
	}

	fg, err := s.finalGradeRepo.GetByUUIDAndUser(ctx, fgUUID, userID)
	if err != nil {
		return fmt.Errorf("get final grade failed: %w", err)
	}
	if fg == nil {
		return fmt.Errorf("final grade not found")
	}

	// Delete associated quiz scores and homework scores
	if err := s.quizScoreRepo.DeleteByFinalGradeID(ctx, fg.ID); err != nil {
		return fmt.Errorf("delete quiz scores failed: %w", err)
	}
	if err := s.homeworkScoreRepo.DeleteByFinalGradeID(ctx, fg.ID); err != nil {
		return fmt.Errorf("delete homework scores failed: %w", err)
	}

	if _, err := s.finalGradeRepo.Delete(ctx, fgUUID, userID); err != nil {
		return fmt.Errorf("delete final grade failed: %w", err)
	}

	return nil
}

// RecalculateFinalScore 重新计算最终成绩
func (s *FinalGradeService) RecalculateFinalScore(ctx context.Context, fg *entity.FinalGrade) error {
	quizAvg, err := s.quizScoreRepo.GetAvgByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return fmt.Errorf("get quiz average failed: %w", err)
	}
	homeworkAvg, err := s.homeworkScoreRepo.GetAvgByFinalGradeID(ctx, fg.ID)
	if err != nil {
		return fmt.Errorf("get homework average failed: %w", err)
	}

	fg.FinalScore = fg.ExamScore*float64(fg.ExamRatio)/100 +
		fg.ClassroomBonusScore*float64(fg.ClassroomBonusRatio)/100 +
		fg.AttendanceScore*float64(fg.AttendanceRatio)/100 +
		quizAvg*float64(fg.QuizRatio)/100 +
		homeworkAvg*float64(fg.HomeworkRatio)/100

	return nil
}

func (s *FinalGradeService) validateRatios(exam, classroomBonus, attendance, quiz, homework int) error {
	total := exam + classroomBonus + attendance + quiz + homework
	if total != 100 {
		return fmt.Errorf("各项比例之和必须为100，当前为%d（考试%d+课堂加分%d+出勤%d+测验%d+作业%d）", total, exam, classroomBonus, attendance, quiz, homework)
	}
	return nil
}
