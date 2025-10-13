package svc

import (
	"time"

	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"
	"shortener/pkg/cache"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	ShortUrlModel        model.ShortUrlMapModel
	ClickStatisticsModel model.ClickStatisticsModel

	Sequence sequence.Sequence

	ShortUrlBlackList map[string]struct{} // 短链接黑名单

	LocalCache *cache.LocalCache // 本地缓存
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 初始化短链接黑名单
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	// 初始化本地缓存，设置TTL为5分钟
	localCache := cache.NewLocalCache(5 * time.Minute)

	return &ServiceContext{
		Config:               c,
		ShortUrlModel:        model.NewShortUrlMapModel(conn, c.CatheRedis),
		ClickStatisticsModel: model.NewClickStatisticsModel(conn, c.CatheRedis),
		// 使用批量Redis序列生成器，批大小设为1000
		Sequence:             sequence.NewBatchRedis(c.CatheRedis[0].Host, 1000),
		ShortUrlBlackList:    m,
		LocalCache:           localCache,
	}
}