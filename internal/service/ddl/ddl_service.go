package ddl

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/schema"
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

type DDLRepo interface {
	AddDraft(ctx context.Context, draft *entity.DDL) error
	GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error)
	GetDraftByUUID(ctx context.Context, uuid string) (*entity.DDL, bool, error)
	UpdateStatusByUUID(ctx context.Context, uuid string, fromStatus int, toStatus int) (int64, error)
	UpdateStatusByUUIDAndUser(ctx context.Context, uuid string, userID int64, fromStatus int, toStatus int) (int64, error)
	GetDDLsForRemind(ctx context.Context, start, end time.Time) ([]*entity.DDL, error)
	MarkRemindSent(ctx context.Context, ddlID int64) error
	GetDDLByID(ctx context.Context, ddlID int64) (*entity.DDL, error)
	GetExpiredDDLs(ctx context.Context, before time.Time) ([]int64, error)
	BatchUpdateStatusToExpired(ctx context.Context, ids []int64) (int64, error)
	DeleteDDLByUUIDAndUser(ctx context.Context, uuid string, userID int64) (int64, error)
	GetDDLsByUserIDAndStatus(ctx context.Context, userID int64, status int, offset, limit int) ([]*entity.DDL, error)
	CountDDLsByUserIDAndStatus(ctx context.Context, userID int64, status int) (int64, error)
	GetDDLByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.DDL, error)
	UpdateDDL(ctx context.Context, ddl *entity.DDL) error
}

type DDLService struct {
	repo       DDLRepo
	aiProvider ai.AIProvider
}

func NewDDLService(repo DDLRepo, aiProvider ai.AIProvider) *DDLService {
	return &DDLService{
		repo:       repo,
		aiProvider: aiProvider,
	}
}

func (s *DDLService) CreateDraft(ctx context.Context, draft *schema.CreateDraftReq, userUUID string) (*schema.CreateDraftResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	draftType := strings.TrimSpace(draft.Type)
	if draftType == "" {
		draftType = schema.DDLTYPEDEFAULT
		draft.Type = draftType
	}

	switch draftType {
	case schema.DDLTYPEDEFAULT:
	case schema.DDLTYPEPICTURE:
		rawBase64 := strings.TrimSpace(draft.RawBase64)
		if rawBase64 == "" {
			return nil, apperrors.ErrPictureDataMissing
		}
		if s.aiProvider == nil {
			return nil, apperrors.ErrAIProviderDisabled
		}

		imageData, err := base64.StdEncoding.DecodeString(rawBase64)
		if err != nil {
			return nil, apperrors.ErrPictureDataInvalid
		}

		imageDraft, err := s.aiProvider.AnalyzeImage(imageData)
		if err != nil {
			return nil, err
		}
		draft.Draft = schema.CreateDraftInput{
			Title:       imageDraft.Title,
			Description: imageDraft.Description,
			Deadline:    imageDraft.Deadline,
			EarlyRemind: imageDraft.EarlyRemind,
		}
	default:
		return nil, fmt.Errorf("unsupported data_type: %s", draftType)
	}

	// 验证截止时间不能早于当前时间
	now := stime.GetCurrentTime()
	if draft.Draft.Deadline.Before(now) {
		return nil, apperrors.ErrDeadlineInPast
	}

	d := &entity.DDL{
		UUID:            uuid.GenerateUUID(),
		UserID:          userID,
		Title:           draft.Draft.Title,
		Description:     draft.Draft.Description,
		DeadLine:        draft.Draft.Deadline,
		EarlyRemindTime: stime.GetTimeBeforeMinutesFrom(draft.Draft.Deadline, draft.Draft.EarlyRemind),
		CreatedAt:       now,
		UpdatedAt:       now,
		Status:          entity.DDLStatusDraft,
	}
	if err := s.repo.AddDraft(ctx, d); err != nil {
		return nil, err
	}
	return &schema.CreateDraftResp{
		UUID:        d.UUID,
		Title:       d.Title,
		Description: d.Description,
		Deadline:    d.DeadLine,
		EarlyRemind: draft.Draft.EarlyRemind,
	}, nil
}

func (s *DDLService) ApproveDraft(ctx context.Context, uuid string, req *schema.UpdateDraftStatusReq, userUUID string) (*schema.UpdateDraftStatusResp, error) {
	targetStatus := strings.TrimSpace(req.Status)
	if targetStatus != schema.DDLSTATUSACTIVE {
		return nil, apperrors.ErrInvalidDraftStatus
	}

	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	// 获取草稿
	draft, exists, err := s.repo.GetDraftByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrDraftNotFound
	}
	if draft.UserID != userID {
		return nil, apperrors.ErrDraftNotOwned
	}

	// 验证截止时间不能早于当前时间
	if draft.DeadLine.Before(stime.GetCurrentTime()) {
		return nil, apperrors.ErrDeadlineInPast
	}

	// 更新状态
	affected, err := s.repo.UpdateStatusByUUIDAndUser(ctx, uuid, userID, entity.DDLStatusDraft, entity.DDLStatusActive)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, apperrors.ErrDraftStateConflict
	}

	return &schema.UpdateDraftStatusResp{
		UUID:   uuid,
		Status: entity.DDLStatusActive,
	}, nil
}

// DeleteDDL 删除 DDL（软删除，更新状态为 deleted）
func (s *DDLService) DeleteDDL(ctx context.Context, uuid string, userUUID string) error {
	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return err
	}

	// 删除 DDL
	affected, err := s.repo.DeleteDDLByUUIDAndUser(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperrors.ErrDDLNotFound
	}
	return nil
}

// GetActiveDDLs 获取用户激活状态的DDL列表
func (s *DDLService) GetActiveDDLs(ctx context.Context, userUUID string, pageReq *schema.PageReq) (*schema.DDLListResp, error) {
	return s.getDDLsByStatus(ctx, userUUID, entity.DDLStatusActive, pageReq)
}

// GetExpiredDDLs 获取用户过期状态的DDL列表
func (s *DDLService) GetExpiredDDLs(ctx context.Context, userUUID string, pageReq *schema.PageReq) (*schema.DDLListResp, error) {
	return s.getDDLsByStatus(ctx, userUUID, entity.DDLStatusExpired, pageReq)
}

func (s *DDLService) getDDLsByStatus(ctx context.Context, userUUID string, status int, pageReq *schema.PageReq) (*schema.DDLListResp, error) {
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	pageReq.Normalize()

	total, err := s.repo.CountDDLsByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	ddls, err := s.repo.GetDDLsByUserIDAndStatus(ctx, userID, status, pageReq.Offset(), pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]schema.DDLListItem, 0, len(ddls))
	for _, d := range ddls {
		list = append(list, schema.DDLListItem{
			UUID:            d.UUID,
			Title:           d.Title,
			Description:     d.Description,
			Deadline:        d.DeadLine,
			EarlyRemindTime: d.EarlyRemindTime,
			Status:          d.Status,
			CreatedAt:       d.CreatedAt,
		})
	}

	return &schema.DDLListResp{
		List:     list,
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}, nil
}

// UpdateDDL 修改DDL
func (s *DDLService) UpdateDDL(ctx context.Context, uuid string, req *schema.UpdateDDLReq, userUUID string) (*schema.UpdateDDLResp, error) {
	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	// 获取DDL
	ddl, err := s.repo.GetDDLByUUIDAndUser(ctx, uuid, userID)
	if err != nil {
		return nil, err
	}
	if ddl == nil {
		return nil, apperrors.ErrDDLNotFound
	}

	// 只允许修改 active 状态的 DDL
	if ddl.Status != entity.DDLStatusActive {
		return nil, apperrors.ErrDDLNotActive
	}

	// 验证截止时间不能早于当前时间
	now := stime.GetCurrentTime()
	effectiveDeadline := ddl.DeadLine
	if req.Deadline != nil {
		effectiveDeadline = *req.Deadline
	}
	if effectiveDeadline.Before(now) {
		return nil, apperrors.ErrDeadlineInPast
	}

	// 计算新的 early_remind_time
	var newEarlyRemindTime time.Time
	var earlyRemind int
	if req.EarlyRemind != nil {
		earlyRemind = *req.EarlyRemind
		if req.Deadline != nil {
			newEarlyRemindTime = stime.GetTimeBeforeMinutesFrom(*req.Deadline, earlyRemind)
		} else {
			newEarlyRemindTime = stime.GetTimeBeforeMinutesFrom(ddl.DeadLine, earlyRemind)
		}
	} else if req.Deadline != nil {
		// deadline 变了但 early_remind 没变，重新计算 early_remind_time
		oldRemindMinutes := int(ddl.DeadLine.Sub(ddl.EarlyRemindTime).Minutes())
		earlyRemind = oldRemindMinutes
		newEarlyRemindTime = stime.GetTimeBeforeMinutesFrom(*req.Deadline, earlyRemind)
	} else {
		earlyRemind = int(ddl.DeadLine.Sub(ddl.EarlyRemindTime).Minutes())
		newEarlyRemindTime = ddl.EarlyRemindTime
	}

	// 更新字段
	needResetRemind := false
	if req.Title != nil {
		ddl.Title = *req.Title
	}
	if req.Description != nil {
		ddl.Description = *req.Description
	}
	if req.Deadline != nil {
		if !ddl.DeadLine.Equal(*req.Deadline) {
			needResetRemind = true
		}
		ddl.DeadLine = *req.Deadline
	}
	if req.EarlyRemind != nil {
		needResetRemind = true
		ddl.EarlyRemindTime = newEarlyRemindTime
	} else if req.Deadline != nil {
		ddl.EarlyRemindTime = newEarlyRemindTime
	}

	// 如果时间相关字段变了，重置提醒状态
	if needResetRemind {
		ddl.RemindSent = false
	}
	ddl.UpdatedAt = now

	// 保存
	if err := s.repo.UpdateDDL(ctx, ddl); err != nil {
		return nil, err
	}

	return &schema.UpdateDDLResp{
		UUID:        ddl.UUID,
		Title:       ddl.Title,
		Description: ddl.Description,
		Deadline:    ddl.DeadLine,
		EarlyRemind: earlyRemind,
	}, nil
}

// GetDDLDetail 获取DDL详情
func (s *DDLService) GetDDLDetail(ctx context.Context, uuid string, userUUID string) (*schema.DDLDetailResp, error) {
	// 获取用户ID
	userID, err := s.repo.GetUserIDByUserUUID(ctx, strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}

	// 获取DDL
	ddl, err := s.repo.GetDDLByUUIDAndUser(ctx, uuid, userID)
	if err != nil {
		return nil, err
	}
	if ddl == nil {
		return nil, apperrors.ErrDDLNotFound
	}

	return &schema.DDLDetailResp{
		UUID:            ddl.UUID,
		Title:           ddl.Title,
		Description:     ddl.Description,
		Deadline:        ddl.DeadLine,
		EarlyRemindTime: ddl.EarlyRemindTime,
		Status:          ddl.Status,
		RemindSent:      ddl.RemindSent,
		CreatedAt:       ddl.CreatedAt,
		UpdatedAt:       ddl.UpdatedAt,
	}, nil
}
