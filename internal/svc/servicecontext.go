package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	ShortUrlModel        model.ShortUrlMapModel
	ClickStatisticsModel model.ClickStatisticsModel

	Sequence sequence.Sequence

	ShortUrlBlackList map[string]struct{} // 短链接黑名单

}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 初始化短链接黑名单
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	return &ServiceContext{
		Config:               c,
		ShortUrlModel:        model.NewShortUrlMapModel(conn, c.CatheRedis),
		ClickStatisticsModel: model.NewClickStatisticsModel(conn, c.CatheRedis),
		Sequence:             sequence.NewMySQL(c.Sequence.DSN),
		//Sequence: sequence.NewRedis(c.CatheRedis[0].Host),
	}
}