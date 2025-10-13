package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"shortener/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
)

// ClickWorker 点击统计异步处理worker
type ClickWorker struct {
	redisClient *redis.Redis
	model       model.ClickStatisticsModel
	queueName   string
	batchSize   int
	quit        chan struct{}
	wg          syncx.WaitGroupWrapper
}

// NewClickWorker 创建点击统计worker
func NewClickWorker(redisClient *redis.Redis, model model.ClickStatisticsModel) *ClickWorker {
	worker := &ClickWorker{
		redisClient: redisClient,
		model:       model,
		queueName:   "shortener:access_queue",
		batchSize:   100,
		quit:        make(chan struct{}),
	}

	return worker
}

// Start 启动worker
func (w *ClickWorker) Start() {
	// 启动多个工作协程
	for i := 0; i < 4; i++ {
		w.wg.Add(func() {
			w.processQueue()
		})
	}

	logx.Info("点击统计worker已启动")
}

// Stop 停止worker
func (w *ClickWorker) Stop() {
	close(w.quit)
	w.wg.Wait()
	logx.Info("点击统计worker已停止")
}

// processQueue 处理队列中的点击记录
func (w *ClickWorker) processQueue() {
	ticker := time.NewTicker(time.Millisecond * 100) // 每100ms尝试从队列获取一次数据
	defer ticker.Stop()

	for {
		select {
		case <-w.quit:
			return
		case <-ticker.C:
			w.batchProcess()
		}
	}
}

// batchProcess 批量处理点击记录
func (w *ClickWorker) batchProcess() {
	// 从队列中批量获取数据
	records, err := w.redisClient.Lrange(w.queueName, 0, int64(w.batchSize-1))
	if err != nil {
		logx.Errorw("从队列获取数据失败", logx.Field{Key: "err", Value: err.Error()})
		return
	}

	if len(records) == 0 {
		return
	}

	// 解析并处理记录
	for _, recordStr := range records {
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(recordStr), &record); err != nil {
			logx.Errorw("解析记录失败", logx.Field{Key: "err", Value: err.Error()})
			continue
		}

		// 保存到数据库
		w.saveToDatabase(record)
	}

	// 从队列中移除已处理的记录
	w.redisClient.Ltrim(w.queueName, int64(w.batchSize), -1)
}

// saveToDatabase 保存点击记录到数据库
func (w *ClickWorker) saveToDatabase(record map[string]interface{}) {
	surl, _ := record["surl"].(string)
	ip, _ := record["ip"].(string)
	userAgent, _ := record["user_agent"].(string)
	referer, _ := record["referer"].(string)

	// 解析时间戳
	clickTime := time.Now()
	if ts, ok := record["click_time"].(float64); ok {
		clickTime = time.Unix(int64(ts), 0)
	}

	// 保存到数据库
	_, err := w.model.Insert(context.Background(), &model.ClickStatistics{
		Surl:      surl,
		ClickTime: clickTime,
		Ip:        sql.NullString{String: ip, Valid: true},
		UserAgent: sql.NullString{String: userAgent, Valid: true},
		Referer:   sql.NullString{String: referer, Valid: referer != ""},
	})

	if err != nil {
		logx.Errorw("保存点击记录失败", logx.Field{Key: "err", Value: err.Error()})
	}
}