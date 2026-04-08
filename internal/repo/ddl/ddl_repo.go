package ddl

import (
	"context"
	"errors"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/data"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
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

func (r *ddlRepo) GetUserIDByUserUUID(ctx context.Context, uuid string) (int64, error) {
	user := &entity.User{UUID: uuid}
	has, err := r.data.DB.Context(ctx).Get(user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, apperrors.ErrUserNotFound
	}
	if user.ID <= 0 {
		return 0, errors.New("invalid user id")
	}

	return user.ID, nil
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

func (r *ddlRepo) UpdateStatusByUUIDAndUser(ctx context.Context, uuid string, userID int64, fromStatus int, toStatus int) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ? AND status = ?", uuid, userID, fromStatus).
		Cols("status", "updated_at").
		Update(&entity.DDL{
			Status:    toStatus,
			UpdatedAt: stime.GetCurrentTime(),
		})
}

// GetDDLsForRemind 获取指定时间范围内需要提醒的 DDL
// 查询条件：status = active，early_remind_time 在 [start, end] 范围内，remind_sent = false
func (r *ddlRepo) GetDDLsForRemind(ctx context.Context, start, end time.Time) ([]*entity.DDL, error) {
	var ddls []*entity.DDL
	err := r.data.DB.Context(ctx).
		Where("status = ?", entity.DDLStatusActive).
		And("early_remind_time >= ?", start).
		And("early_remind_time <= ?", end).
		And("remind_sent = ?", false).
		And("early_remind_time IS NOT NULL").
		Find(&ddls)
	if err != nil {
		return nil, err
	}
	return ddls, nil
}

// MarkRemindSent 标记提醒已发送
func (r *ddlRepo) MarkRemindSent(ctx context.Context, ddlID int64) error {
	_, err := r.data.DB.Context(ctx).
		ID(ddlID).
		Cols("remind_sent", "updated_at").
		Update(&entity.DDL{RemindSent: true})
	return err
}

// GetDDLByID 根据 ID 获取 DDL
func (r *ddlRepo) GetDDLByID(ctx context.Context, ddlID int64) (*entity.DDL, error) {
	ddl := &entity.DDL{}
	has, err := r.data.DB.Context(ctx).ID(ddlID).Get(ddl)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return ddl, nil
}

// GetDDLsForRemindWithUserEmail 获取需要提醒的 DDL 及用户邮箱
func (r *ddlRepo) GetDDLsForRemindWithUserEmail(ctx context.Context, start, end time.Time) ([]*entity.DDLWithUserEmail, error) {
	var results []*entity.DDLWithUserEmail
	err := r.data.DB.Context(ctx).
		Table("ddl").
		Join("INNER", "user", "ddl.user_id = user.id").
		Where("ddl.status = ?", entity.DDLStatusActive).
		And("ddl.early_remind_time >= ?", start).
		And("ddl.early_remind_time <= ?", end).
		And("ddl.remind_sent = ?", false).
		And("ddl.early_remind_time IS NOT NULL").
		And("user.email IS NOT NULL").
		And("user.email != ''").
		Cols("ddl.*", "user.email").
		Find(&results)
	return results, err
}
