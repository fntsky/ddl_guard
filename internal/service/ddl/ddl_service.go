package ddl

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

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
			EealyRemind: imageDraft.EealyRemind,
		}
	default:
		return nil, fmt.Errorf("unsupported data_type: %s", draftType)
	}

	d := &entity.DDL{
		UUID:            uuid.GenerateUUID(),
		UserID:          userID,
		Title:           draft.Draft.Title,
		Description:     draft.Draft.Description,
		DeadLine:        draft.Draft.Deadline,
		EealyRemindTime: stime.GetTimeBeforeMinutes(draft.Draft.EealyRemind),
		CreatedAt:       stime.GetCurrentTime(),
		UpdatedAt:       stime.GetCurrentTime(),
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
		EealyRemind: draft.Draft.EealyRemind,
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

	// 只更新属于该用户的草稿
	affected, err := s.repo.UpdateStatusByUUIDAndUser(ctx, uuid, userID, entity.DDLStatusDraft, entity.DDLStatusActive)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		_, exists, err := s.repo.GetDraftByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apperrors.ErrDraftNotFound
		}
		// 草稿存在但不属于该用户
		return nil, apperrors.ErrDraftNotOwned
	}

	return &schema.UpdateDraftStatusResp{
		UUID:   uuid,
		Status: entity.DDLStatusActive,
	}, nil
}
