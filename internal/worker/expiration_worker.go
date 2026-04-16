package worker

import (
	"context"
	"log"
	"sync"
	"time"

	ddlsvc "github.com/fntsky/ddl_guard/internal/service/ddl"
	stime "github.com/fntsky/ddl_guard/pkg/time"
)

const (
	expirationScanInterval = 1 * time.Minute // 扫描间隔
	expirationBatchSize    = 100             // 每批处理数量
)

type ExpirationWorker struct {
	ddlRepo ddlsvc.DDLRepo
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewExpirationWorker(ddlRepo ddlsvc.DDLRepo) *ExpirationWorker {
	return &ExpirationWorker{
		ddlRepo: ddlRepo,
		stopCh:  make(chan struct{}),
	}
}

func (w *ExpirationWorker) Start() {
	w.wg.Add(1)
	go w.run()
	log.Println("[ExpirationWorker] started")
}

func (w *ExpirationWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	log.Println("[ExpirationWorker] stopped")
}

func (w *ExpirationWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(expirationScanInterval)
	defer ticker.Stop()

	// 启动时立即执行一次
	w.checkExpiredDDLs()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.checkExpiredDDLs()
		}
	}
}

func (w *ExpirationWorker) checkExpiredDDLs() {
	ctx := context.Background()
	now := stime.GetCurrentTime()

	// 获取已过期的 DDL ID
	ids, err := w.ddlRepo.GetExpiredDDLs(ctx, now)
	if err != nil {
		log.Printf("[ExpirationWorker] failed to get expired DDLs: %v", err)
		return
	}

	if len(ids) == 0 {
		return
	}

	// 分批更新
	for i := 0; i < len(ids); i += expirationBatchSize {
		end := i + expirationBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]

		affected, err := w.ddlRepo.BatchUpdateStatusToExpired(ctx, batch)
		if err != nil {
			log.Printf("[ExpirationWorker] failed to update batch: %v", err)
			continue
		}
		log.Printf("[ExpirationWorker] updated %d DDLs to expired status", affected)
	}
}
