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
	ddlsvc "github.com/fntsky/ddl_guard/internal/service/ddl"
	usersvc "github.com/fntsky/ddl_guard/internal/service/user"
	stime "github.com/fntsky/ddl_guard/pkg/time"
	"github.com/redis/go-redis/v9"
)

const (
	scanInterval    = 1 * time.Minute      // 扫描间隔
	preloadMinutes  = 10                   // 预加载时间（分钟）
	remindZSetKey   = "ddl:remind:pending" // ZSET: 排序用
	remindDetailKey = "ddl:remind:detail"  // Hash: 存储详情
)

// DDLCache 用于存储在 Redis 中的 DDL 信息
type DDLCache struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    int64  `json:"deadline"` // Unix timestamp
	Email       string `json:"email"`    // 用户邮箱
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
func (w *PublishWorker) preloadToRedis(ctx context.Context, now time.Time) {
	if w.redis == nil || w.redis.Client == nil {
		return
	}

	// 允许获取过去30分钟内的提醒（处理启动时错过的）
	pastStart := now.Add(-30 * time.Minute)
	end := now.Add(preloadMinutes * time.Minute)

	// 1. 查询 DDL
	ddls, err := w.ddlRepo.GetDDLsForRemind(ctx, pastStart, end)
	if err != nil {
		log.Printf("[PublishWorker] failed to get DDLs: %v", err)
		return
	}
	if len(ddls) == 0 {
		return
	}

	// 2. 收集用户ID，批量查询邮箱
	userIDs := make([]int64, 0, len(ddls))
	for _, d := range ddls {
		userIDs = append(userIDs, d.UserID)
	}
	userEmailMap, err := w.userRepo.GetUserEmailsByIDs(ctx, userIDs)
	if err != nil {
		log.Printf("[PublishWorker] failed to get user emails: %v", err)
		return
	}

	// 3. 写入 Redis
	for _, d := range ddls {
		email, ok := userEmailMap[d.UserID]
		if !ok {
			continue // 用户没有邮箱
		}

		ddlIDStr := fmt.Sprintf("%d", d.ID)

		// 检查是否已在 Redis 中
		score, err := w.redis.Client.ZScore(ctx, remindZSetKey, ddlIDStr).Result()
		if err != nil && err != redis.Nil {
			log.Printf("[PublishWorker] failed to check redis: %v", err)
			continue
		}
		if score != 0 {
			continue // 已存在
		}

		// 使用 Pipeline 批量写入
		pipe := w.redis.Client.Pipeline()

		// 1. 加入 ZSET
		pipe.ZAdd(ctx, remindZSetKey, redis.Z{
			Score:  float64(d.EarlyRemindTime.Unix()),
			Member: d.ID,
		})

		// 2. 存储详情到 Hash
		cache := DDLCache{
			ID:          d.ID,
			UserID:      d.UserID,
			Title:       d.Title,
			Description: d.Description,
			Deadline:    d.DeadLine.Unix(),
			Email:       email,
		}
		data, _ := json.Marshal(cache)
		pipe.HSet(ctx, remindDetailKey, ddlIDStr, data)

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

	nowTimestamp := now.Unix()
	ddlIDs, err := w.redis.Client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     remindZSetKey,
		Start:   "-inf",
		Stop:    fmt.Sprintf("%d", nowTimestamp),
		ByScore: true,
	}).Result()
	if err != nil {
		log.Printf("[PublishWorker] failed to get from redis: %v", err)
		return
	}

	for _, ddlIDStr := range ddlIDs {
		// 从 Hash 获取详情
		data, err := w.redis.Client.HGet(ctx, remindDetailKey, ddlIDStr).Result()
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

		// 直接发送，无需查数据库
		if err := w.sendEmailNotificationFromCache(ctx, &cache); err != nil {
			log.Printf("[PublishWorker] failed to send notification for DDL %d: %v", cache.ID, err)
		} else {
			// 清理：移除 ZSET 和 Hash
			w.redis.Client.ZRem(ctx, remindZSetKey, ddlIDStr)
			w.redis.Client.HDel(ctx, remindDetailKey, ddlIDStr)
			w.ddlRepo.MarkRemindSent(ctx, cache.ID)
			log.Printf("[PublishWorker] notification for DDL %d sent successfully", cache.ID)
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
