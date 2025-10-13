package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 1. 首先尝试从本地缓存获取
	if longUrl, ok := l.svcCtx.LocalCache.Get(req.ShortUrl); ok {
		// 异步记录访问
		go l.recordAccessWithRedisQueue(req.ShortUrl)
		return &types.ShowResponse{LongUrl: longUrl},
			nil
	}

	// 2. 从数据库/Redis获取长链接
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{Valid: true, String: req.ShortUrl})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("404")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Value: err.Error(), Key: "err"})
		return nil, err
	}

	// 3. 将结果存入本地缓存
	l.svcCtx.LocalCache.Set(req.ShortUrl, u.Lurl.String)

	// 4. 异步处理点击记录，使用Redis队列代替直接数据库操作
	go l.recordAccessWithRedisQueue(req.ShortUrl)

	// 5. 返回长链接
	return &types.ShowResponse{
			LongUrl: u.Lurl.String,
		},
		nil
}

// recordAccessWithRedisQueue 使用Redis队列记录访问，而不是直接写数据库
func (l *ShowLogic) recordAccessWithRedisQueue(shortUrl string) {
	// 获取请求信息
	ip := ""
	if v, ok := l.ctx.Value("ip").(string); ok {
		ip = v
	}
	userAgent := ""
	if v, ok := l.ctx.Value("userAgent").(string); ok {
		userAgent = v
	}
	referer := ""
	if v, ok := l.ctx.Value("referer").(string); ok {
		referer = v
	}

	// 构建访问记录
	accessRecord := map[string]interface{}{
		"surl":       shortUrl,
		"click_time": time.Now().Unix(),
		"ip":         ip,
		"user_agent": userAgent,
		"referer":    referer,
	}

	// 序列化为JSON并推入Redis队列
	data, _ := json.Marshal(accessRecord)

	// 创建Redis客户端
	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: l.svcCtx.Config.CatheRedis[0].Host,
	})

	_, err := redisClient.Lpush("shortener:access_queue", string(data))
	if err != nil {
		logx.Errorw("推入访问记录队列失败", logx.LogField{Key: "err", Value: err.Error()})
		// 如果队列失败，降级到直接写库
		l.recordAccessDirectly(shortUrl, ip, userAgent, referer)
	}
}

// recordAccessDirectly 直接写入数据库（降级方案）
func (l *ShowLogic) recordAccessDirectly(shortUrl, ip, userAgent, referer string) {
	_, err := l.svcCtx.ClickStatisticsModel.Insert(l.ctx, &model.ClickStatistics{
		Surl:      shortUrl,
		ClickTime: time.Now(),
		Ip:        sql.NullString{String: ip, Valid: true},
		UserAgent: sql.NullString{String: userAgent, Valid: true},
		Referer:   sql.NullString{String: referer, Valid: referer != ""},
	})
	if err != nil {
		logx.Errorw("保存点击记录失败", logx.LogField{Key: "err", Value: err.Error()})
	}
}
