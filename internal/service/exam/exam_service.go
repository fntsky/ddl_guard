package exam

import (
	"context"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/repo/exam"
	"github.com/fntsky/ddl_guard/internal/schema"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

// ExamService 考试服务
type ExamService struct {
	repo exam.ExamRepo
}

func NewExamService(repo exam.ExamRepo) *ExamService {
	return &ExamService{repo: repo}
}

// CreateExam 创建考试
func (s *ExamService) CreateExam(ctx context.Context, req *schema.CreateExamReq, userUUID string) (*schema.CreateExamResp, error) {
	// 验证时间：结束时间必须晚于开始时间
	if !req.EndTime.After(req.StartTime) {
		return nil, apperrors.ErrExamTimeInvalid
	}

	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	now := stime.GetCurrentTime()
	e := &entity.Exam{
		UUID:      uuid.GenerateUUID(),
		UserID:    userID,
		Name:      strings.TrimSpace(req.Name),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Location:  strings.TrimSpace(req.Location),
		Notes:     strings.TrimSpace(req.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}

	return &schema.CreateExamResp{
		UUID:      e.UUID,
		Name:      e.Name,
		StartTime: e.StartTime,
		EndTime:   e.EndTime,
		Location:  e.Location,
		Notes:     e.Notes,
	}, nil
}

// GetExam 获取考试详情
func (s *ExamService) GetExam(ctx context.Context, examUUID string, userUUID string) (*schema.ExamDetailResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	exam, err := s.repo.GetByUUIDAndUser(ctx, examUUID, userID)
	if err != nil {
		return nil, err
	}
	if exam == nil {
		return nil, apperrors.ErrExamNotFound
	}

	return &schema.ExamDetailResp{
		UUID:      exam.UUID,
		Name:      exam.Name,
		StartTime: exam.StartTime,
		EndTime:   exam.EndTime,
		Location:  exam.Location,
		Notes:     exam.Notes,
		CreatedAt: exam.CreatedAt,
		UpdatedAt: exam.UpdatedAt,
	}, nil
}

// ListExams 获取考试列表
func (s *ExamService) ListExams(ctx context.Context, userUUID string, pageReq *schema.PageReq) (*schema.ExamListResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	pageReq.Normalize()

	total, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	exams, err := s.repo.ListByUserID(ctx, userID, pageReq.Offset(), pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]schema.ExamListItem, 0, len(exams))
	for _, e := range exams {
		list = append(list, schema.ExamListItem{
			UUID:      e.UUID,
			Name:      e.Name,
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
			Location:  e.Location,
			CreatedAt: e.CreatedAt,
		})
	}

	return &schema.ExamListResp{
		List:     list,
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

// UpdateExam 更新考试
func (s *ExamService) UpdateExam(ctx context.Context, examUUID string, req *schema.UpdateExamReq, userUUID string) (*schema.UpdateExamResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	exam, err := s.repo.GetByUUIDAndUser(ctx, examUUID, userID)
	if err != nil {
		return nil, err
	}
	if exam == nil {
		return nil, apperrors.ErrExamNotFound
	}

	// 更新字段
	if req.Name != nil {
		exam.Name = strings.TrimSpace(*req.Name)
	}
	if req.StartTime != nil {
		exam.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		exam.EndTime = *req.EndTime
	}
	if req.Location != nil {
		exam.Location = strings.TrimSpace(*req.Location)
	}
	if req.Notes != nil {
		exam.Notes = strings.TrimSpace(*req.Notes)
	}

	// 验证时间：结束时间必须晚于开始时间
	if !exam.EndTime.After(exam.StartTime) {
		return nil, apperrors.ErrExamTimeInvalid
	}

	exam.UpdatedAt = stime.GetCurrentTime()

	if err := s.repo.Update(ctx, exam); err != nil {
		return nil, err
	}

	return &schema.UpdateExamResp{
		UUID:      exam.UUID,
		Name:      exam.Name,
		StartTime: exam.StartTime,
		EndTime:   exam.EndTime,
		Location:  exam.Location,
		Notes:     exam.Notes,
	}, nil
}

// DeleteExam 删除考试
func (s *ExamService) DeleteExam(ctx context.Context, examUUID string, userUUID string) error {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return err
	}

	affected, err := s.repo.Delete(ctx, examUUID, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperrors.ErrExamNotFound
	}

	return nil
}
