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

// GetExpiredDDLs 获取已过期的 DDL ID 列表
// 查询条件：status = active，deadline < before
func (r *ddlRepo) GetExpiredDDLs(ctx context.Context, before time.Time) ([]int64, error) {
	var ids []int64
	err := r.data.DB.Context(ctx).
		Table("ddl").
		Where("status = ?", entity.DDLStatusActive).
		And("deadline < ?", before).
		Cols("id").
		Find(&ids)
	return ids, err
}

// BatchUpdateStatusToExpired 批量更新 DDL 状态为过期
func (r *ddlRepo) BatchUpdateStatusToExpired(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return r.data.DB.Context(ctx).
		In("id", ids).
		Cols("status", "updated_at").
		Update(&entity.DDL{
			Status:    entity.DDLStatusExpired,
			UpdatedAt: stime.GetCurrentTime(),
		})
}

// DeleteDDLByUUIDAndUser 删除 DDL（软删除，更新状态为 deleted）
func (r *ddlRepo) DeleteDDLByUUIDAndUser(ctx context.Context, uuid string, userID int64) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ? AND status != ?", uuid, userID, entity.DDLStatusDeleted).
		Cols("status", "updated_at").
		Update(&entity.DDL{
			Status:    entity.DDLStatusDeleted,
			UpdatedAt: stime.GetCurrentTime(),
		})
}

// GetDDLsByUserIDAndStatus 分页查询用户的DDL
func (r *ddlRepo) GetDDLsByUserIDAndStatus(ctx context.Context, userID int64, status int, offset, limit int) ([]*entity.DDL, error) {
	var ddls []*entity.DDL
	err := r.data.DB.Context(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Desc("deadline").
		Limit(limit, offset).
		Find(&ddls)
	return ddls, err
}

// CountDDLsByUserIDAndStatus 统计用户DDL数量
func (r *ddlRepo) CountDDLsByUserIDAndStatus(ctx context.Context, userID int64, status int) (int64, error) {
	return r.data.DB.Context(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Count(&entity.DDL{})
}

// GetDDLByUUIDAndUser 获取用户的DDL
func (r *ddlRepo) GetDDLByUUIDAndUser(ctx context.Context, uuid string, userID int64) (*entity.DDL, error) {
	ddl := &entity.DDL{}
	has, err := r.data.DB.Context(ctx).
		Where("uuid = ? AND user_id = ?", uuid, userID).
		Get(ddl)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return ddl, nil
}

// UpdateDDL 更新DDL
func (r *ddlRepo) UpdateDDL(ctx context.Context, ddl *entity.DDL) error {
	_, err := r.data.DB.Context(ctx).
		ID(ddl.ID).
		AllCols().
		Update(ddl)
	return err
}
