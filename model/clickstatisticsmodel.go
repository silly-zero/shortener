package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ClickStatisticsModel = (*customClickStatisticsModel)(nil)

type (
	// ClickStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClickStatisticsModel.
	ClickStatisticsModel interface {
		clickStatisticsModel
		// FindBySurlAndDateRange 根据短链接和日期范围查询点击记录
		FindBySurlAndDateRange(ctx context.Context, surl string, startDate, endDate time.Time) ([]*ClickStatistics, error)
		// CountTotalClicks 统计短链接的总点击量
		CountTotalClicks(ctx context.Context, surl string, startDate, endDate time.Time) (int64, error)
		// CountTotalUniqueVisitors 统计短链接的总唯一访问者数量
		CountTotalUniqueVisitors(ctx context.Context, surl string, startDate, endDate time.Time) (int64, error)
		// GetDailyStats 获取短链接的每日统计数据
		GetDailyStats(ctx context.Context, surl string, startDate, endDate time.Time) ([]DailyStatsRow, error)
	}

	customClickStatisticsModel struct {
		*defaultClickStatisticsModel
	}

	// DailyStatsRow 每日统计数据行
	DailyStatsRow struct {
		ClickDate      string
		TotalClicks    int64
		UniqueVisitors int64
	}
)

// NewClickStatisticsModel returns a model for the database table.
func NewClickStatisticsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ClickStatisticsModel {
	return &customClickStatisticsModel{
		defaultClickStatisticsModel: newClickStatisticsModel(conn, c, opts...),
	}
}

// FindBySurlAndDateRange 根据短链接和日期范围查询点击记录
func (m *customClickStatisticsModel) FindBySurlAndDateRange(ctx context.Context, surl string, startDate, endDate time.Time) ([]*ClickStatistics, error) {
	var clickStats []*ClickStatistics
	key := "custom:clickstatistics:findbysurl:range"

	// 使用ExecCtx包装查询逻辑
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		// 准备SQL语句
		query := `select id, surl, click_time, ip, user_agent, referer from click_statistics where surl = ? and click_time between ? and ?`

		// 使用stmtSession来执行多行查询
		stmt, err := conn.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		// 使用QueryRowsCtx直接执行多行查询
		if err := stmt.QueryRowsCtx(ctx, &clickStats, surl, startDate, endDate); err != nil {
			// 如果没有记录，返回空切片
			if err == sqlx.ErrNotFound {
				return nil, nil
			}
			return nil, err
		}

		return nil, nil
	}, key)

	return clickStats, err
}

// CountTotalClicks 统计短链接的总点击量
func (m *customClickStatisticsModel) CountTotalClicks(ctx context.Context, surl string, startDate, endDate time.Time) (int64, error) {
	var count int64
	key := "custom:clickstatistics:count:total"

	// 使用QueryRowCtx执行计数查询
	err := m.QueryRowCtx(ctx, &count, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := `select count(*) from click_statistics where surl = ? and click_time between ? and ?`
		return conn.QueryRowCtx(ctx, v, query, surl, startDate, endDate)
	})

	return count, err
}

// CountTotalUniqueVisitors 统计短链接的总唯一访问者数量
func (m *customClickStatisticsModel) CountTotalUniqueVisitors(ctx context.Context, surl string, startDate, endDate time.Time) (int64, error) {
	var count int64
	key := "custom:clickstatistics:count:unique"

	// 使用QueryRowCtx执行唯一计数查询
	err := m.QueryRowCtx(ctx, &count, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := `select count(distinct ip) from click_statistics where surl = ? and click_time between ? and ?`
		return conn.QueryRowCtx(ctx, v, query, surl, startDate, endDate)
	})

	return count, err
}

// GetDailyStats 获取短链接的每日统计数据
func (m *customClickStatisticsModel) GetDailyStats(ctx context.Context, surl string, startDate, endDate time.Time) ([]DailyStatsRow, error) {
	var dailyStats []DailyStatsRow
	key := "custom:clickstatistics:daily:stats"

	// 使用ExecCtx包装查询逻辑
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		// 直接使用数据库中的视图来简化查询
		query := `select click_date, total_clicks, unique_visitors from daily_click_summary where surl = ? and click_date between date(?) and date(?) order by click_date`

		// 使用stmtSession来执行多行查询
		stmt, err := conn.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		// 使用QueryRowsCtx直接执行多行查询
		if err := stmt.QueryRowsCtx(ctx, &dailyStats, surl, startDate, endDate); err != nil {
			// 如果没有记录，返回空切片
			if err == sqlx.ErrNotFound {
				return nil, nil
			}
			return nil, err
		}

		return nil, nil
	}, key)

	return dailyStats, err
}