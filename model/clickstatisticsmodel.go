package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ClickStatisticsModel = (*customClickStatisticsModel)(nil)

type (
	// ClickStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClickStatisticsModel.
	ClickStatisticsModel interface {
		clickStatisticsModel
	}

	customClickStatisticsModel struct {
		*defaultClickStatisticsModel
	}
)

// NewClickStatisticsModel returns a model for the database table.
func NewClickStatisticsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ClickStatisticsModel {
	return &customClickStatisticsModel{
		defaultClickStatisticsModel: newClickStatisticsModel(conn, c, opts...),
	}
}
