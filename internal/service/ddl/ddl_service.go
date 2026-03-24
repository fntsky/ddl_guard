package ddl

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/schema"
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

type DDLRepo interface {
	AddDraft(ctx context.Context, draft *entity.DDL) error
	GetDraftByUUID(ctx context.Context, uuid string) (*entity.DDL, bool, error)
	UpdateStatusByUUID(ctx context.Context, uuid string, fromStatus int, toStatus int) (int64, error)
}

type DDLService struct {
	repo       DDLRepo
	aiProvider ai.AIProvider
}

var (
	ErrInvalidDraftStatus = errors.New("invalid draft status")
	ErrDraftNotFound      = errors.New("draft not found")
	ErrDraftStateConflict = errors.New("draft state conflict")
	ErrPictureDataMissing = errors.New("picture base64 data is required")
	ErrPictureDataInvalid = errors.New("invalid picture base64 data")
	ErrAIProviderDisabled = errors.New("ai provider is not configured")
)

func NewDDLService(repo DDLRepo, aiProvider ai.AIProvider) *DDLService {
	return &DDLService{
		repo:       repo,
		aiProvider: aiProvider,
	}
}

func (s *DDLService) CreateDraft(ctx context.Context, draft *schema.CreateDraftReq) (*schema.CreateDraftResp, error) {
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
			return nil, ErrPictureDataMissing
		}
		if s.aiProvider == nil {
			return nil, ErrAIProviderDisabled
		}

		imageData, err := base64.StdEncoding.DecodeString(rawBase64)
		if err != nil {
			return nil, ErrPictureDataInvalid
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

func (s *DDLService) ApproveDraft(ctx context.Context, uuid string, req *schema.UpdateDraftStatusReq) (*schema.UpdateDraftStatusResp, error) {
	targetStatus := strings.TrimSpace(req.Status)
	if targetStatus != schema.DDLSTATUSACTIVE {
		return nil, ErrInvalidDraftStatus
	}

	affected, err := s.repo.UpdateStatusByUUID(ctx, uuid, entity.DDLStatusDraft, entity.DDLStatusActive)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		_, exists, err := s.repo.GetDraftByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrDraftNotFound
		}
		return nil, ErrDraftStateConflict
	}

	return &schema.UpdateDraftStatusResp{
		UUID:   uuid,
		Status: entity.DDLStatusActive,
	}, nil
}
