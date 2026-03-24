package ddl

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	stime "github.com/fntsky/ddl_guard/pkg/time"
)

type ddlRepo struct {
	data *data.Data
}

func NewDDLRepo(data *data.Data) ddl.DDLRepo {
	return &ddlRepo{
		data: data,
	}
}

// AddDraft 添加一个新的DDL草稿
func (r *ddlRepo) AddDraft(ctx context.Context, draft *entity.DDL) error {
	_, err := r.data.DB.Context(ctx).Insert(draft)
	return err
}

func (r *ddlRepo) GetDraftByUUID(ctx context.Context, uuid string) (*entity.DDL, bool, error) {
	draft := &entity.DDL{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(draft)
	return draft, has, err
}

func (r *ddlRepo) UpdateStatusByUUID(ctx context.Context, uuid string, fromStatus int, toStatus int) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND status = ?", uuid, fromStatus).
		Cols("status", "updated_at").
		Update(&entity.DDL{
			Status:    toStatus,
			UpdatedAt: stime.GetCurrentTime(),
		})
}
