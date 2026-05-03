package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/base/email"
	"github.com/fntsky/ddl_guard/internal/entity"
	ddlsvc "github.com/fntsky/ddl_guard/internal/service/ddl"
	usersvc "github.com/fntsky/ddl_guard/internal/service/user"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/redis/go-redis/v9"
)

const (
	scanInterval      = 1 * time.Minute       // 扫描间隔
	preloadMinutes    = 10                    // 预加载时间（分钟）
	remindZSetKey24h  = "ddl:remind:24h"      // ZSET: 24小时提醒
	remindZSetKey2h   = "ddl:remind:2h"       // ZSET: 2小时提醒
	remindDetailKey   = "ddl:remind:detail"   // Hash: 存储详情
)

// DDLCache 用于存储在 Redis 中的 DDL 信息
type DDLCache struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    int64  `json:"deadline"` // Unix timestamp
	Email       string `json:"email"`    // 用户邮箱
	RemindType  string `json:"remind_type"` // "24h" 或 "2h"
}

type PublishWorker struct {
	ddlRepo     ddlsvc.DDLRepo
	userRepo    usersvc.UserRepo
	emailSender email.Sender
	redis       *data.RedisClient
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func NewPublishWorker(
	ddlRepo ddlsvc.DDLRepo,
	userRepo usersvc.UserRepo,
	redis *data.RedisClient,
) *PublishWorker {
	var sender email.Sender
	cfg := conf.Global()
	if cfg != nil && cfg.Publish.Email.Enabled {
		smtpCfg := cfg.Publish.Email.SMTP
		sender = email.NewSMTPSender(smtpCfg.Host, smtpCfg.Port, smtpCfg.Username, smtpCfg.Password)
	}
	return &PublishWorker{
		ddlRepo:     ddlRepo,
		userRepo:    userRepo,
		emailSender: sender,
		redis:       redis,
		stopCh:      make(chan struct{}),
	}
}

func (w *PublishWorker) Start() {
	w.wg.Add(1)
	go w.run()
	log.Println("[PublishWorker] started")
}

func (w *PublishWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	log.Println("[PublishWorker] stopped")
}

func (w *PublishWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	// 启动时立即执行一次
	w.scanAndNotify()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.scanAndNotify()
		}
	}
}

func (w *PublishWorker) scanAndNotify() {
	ctx := context.Background()
	now := stime.GetCurrentTime()

	// 1. 预加载：扫描数据库中即将到期的 DDL，加入 Redis
	w.preloadToRedis(ctx, now)

	// 2. 从 Redis 获取到期的通知并发送
	w.sendPendingNotifications(ctx, now)
}

// preloadToRedis 将即将到期的 DDL 加入 Redis (ZSET + Hash)
// 提醒规则：截止时间前24小时和前2小时各提醒一次
func (w *PublishWorker) preloadToRedis(ctx context.Context, now time.Time) {
	if w.redis == nil || w.redis.Client == nil {
		return
	}

	// 预加载24小时提醒：deadline 在 [now+24h, now+24h+preloadInterval]
	start24h := now.Add(24 * time.Hour)
	end24h := now.Add(24*time.Hour + preloadMinutes*time.Minute)

	ddls24h, err := w.ddlRepo.GetDDLsForRemind24h(ctx, start24h, end24h)
	if err != nil {
		log.Printf("[PublishWorker] failed to get DDLs for 24h remind: %v", err)
	} else {
		w.preloadDDLsToRedis(ctx, ddls24h, "24h", now)
	}

	// 预加载2小时提醒：deadline 在 [now+2h, now+2h+preloadInterval]
	start2h := now.Add(2 * time.Hour)
	end2h := now.Add(2*time.Hour + preloadMinutes*time.Minute)

	ddls2h, err := w.ddlRepo.GetDDLsForRemind2h(ctx, start2h, end2h)
	if err != nil {
		log.Printf("[PublishWorker] failed to get DDLs for 2h remind: %v", err)
	} else {
		w.preloadDDLsToRedis(ctx, ddls2h, "2h", now)
	}
}

func (w *PublishWorker) preloadDDLsToRedis(ctx context.Context, ddls []*entity.DDL, remindType string, now time.Time) {
	if len(ddls) == 0 {
		return
	}

	userIDs := make([]int64, 0, len(ddls))
	for _, d := range ddls {
		userIDs = append(userIDs, d.UserID)
	}
	userEmailMap, err := w.userRepo.GetUserEmailsByIDs(ctx, userIDs)
	if err != nil {
		log.Printf("[PublishWorker] failed to get user emails: %v", err)
		return
	}

	zsetKey := remindZSetKey24h
	if remindType == "2h" {
		zsetKey = remindZSetKey2h
	}

	for _, d := range ddls {
		email, ok := userEmailMap[d.UserID]
		if !ok {
			continue
		}

		ddlIDStr := fmt.Sprintf("%d", d.ID)

		score, err := w.redis.Client.ZScore(ctx, zsetKey, ddlIDStr).Result()
		if err != nil && err != redis.Nil {
			log.Printf("[PublishWorker] failed to check redis: %v", err)
			continue
		}
		if score != 0 {
			continue
		}

		pipe := w.redis.Client.Pipeline()

		pipe.ZAdd(ctx, zsetKey, redis.Z{
			Score:  float64(d.DeadLine.Unix()),
			Member: d.ID,
		})

		cache := DDLCache{
			ID:          d.ID,
			UserID:      d.UserID,
			Title:       d.Title,
			Description: d.Description,
			Deadline:    d.DeadLine.Unix(),
			Email:       email,
			RemindType:  remindType,
		}
		data, _ := json.Marshal(cache)
		pipe.HSet(ctx, remindDetailKey, ddlIDStr+"_"+remindType, data)

		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("[PublishWorker] failed to add to redis: %v", err)
		}
	}
}

// sendPendingNotifications 从 Redis 获取到期的通知并发送
func (w *PublishWorker) sendPendingNotifications(ctx context.Context, now time.Time) {
	if w.redis == nil || w.redis.Client == nil {
		return
	}

	// 处理24小时提醒（deadline <= now + 24h，即提前24小时已过）
	w.sendReminders(ctx, now, remindZSetKey24h, "24h")

	// 处理2小时提醒
	w.sendReminders(ctx, now, remindZSetKey2h, "2h")
}

func (w *PublishWorker) sendReminders(ctx context.Context, now time.Time, zsetKey, remindType string) {
	// 查找需要发送的提醒：deadline <= now + remindTime
	// 对于24h提醒：当 now >= deadline - 24h 时发送，即 deadline <= now + 24h
	// 对于2h提醒：当 now >= deadline - 2h 时发送，即 deadline <= now + 2h
	var remindOffset time.Duration
	if remindType == "24h" {
		remindOffset = 24 * time.Hour
	} else {
		remindOffset = 2 * time.Hour
	}

	threshold := now.Add(remindOffset).Unix()
	ddlIDs, err := w.redis.Client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     zsetKey,
		Start:   "-inf",
		Stop:    fmt.Sprintf("%d", threshold),
		ByScore: true,
	}).Result()
	if err != nil {
		log.Printf("[PublishWorker] failed to get from redis: %v", err)
		return
	}

	for _, ddlIDStr := range ddlIDs {
		data, err := w.redis.Client.HGet(ctx, remindDetailKey, ddlIDStr+"_"+remindType).Result()
		if err != nil {
			if err != redis.Nil {
				log.Printf("[PublishWorker] failed to get detail: %v", err)
			}
			continue
		}

		var cache DDLCache
		if err := json.Unmarshal([]byte(data), &cache); err != nil {
			log.Printf("[PublishWorker] failed to unmarshal: %v", err)
			continue
		}

		if err := w.sendEmailNotificationFromCache(ctx, &cache); err != nil {
			log.Printf("[PublishWorker] failed to send notification for DDL %d: %v", cache.ID, err)
		} else {
			w.redis.Client.ZRem(ctx, zsetKey, ddlIDStr)
			w.redis.Client.HDel(ctx, remindDetailKey, ddlIDStr+"_"+remindType)
			if remindType == "24h" {
				w.ddlRepo.MarkRemind24hSent(ctx, cache.ID)
			} else {
				w.ddlRepo.MarkRemind2hSent(ctx, cache.ID)
			}
			log.Printf("[PublishWorker] %s notification for DDL %d sent successfully", remindType, cache.ID)
		}
	}
}

func (w *PublishWorker) sendEmailNotificationFromCache(ctx context.Context, cache *DDLCache) error {
	if w.emailSender == nil {
		return fmt.Errorf("email sender not configured")
	}

	subject := fmt.Sprintf("[DDL提醒] %s", cache.Title)
	body := fmt.Sprintf(
		"您好，\n\n您的 DDL 即将到期：\n\n标题：%s\n描述：%s\n截止时间：%s\n\n请及时处理。",
		cache.Title,
		cache.Description,
		time.Unix(cache.Deadline, 0).Format("2006-01-02 15:04"),
	)

	return w.emailSender.Send(ctx, cache.Email, subject, body)
}
