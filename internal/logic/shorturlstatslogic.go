package logic

import (
	"context"
	"database/sql"
	"errors"
	"shortener/model"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortUrlStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortUrlStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortUrlStatsLogic {
	return &ShortUrlStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortUrlStatsLogic) ShortUrlStats(req *types.ShortUrlStatsRequest) (resp *types.ShortUrlStatsResponse, err error) {
	//1.验证短链接是否存在
	if _, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true}); err != nil {
		return nil, errors.New("short url not found")
	}

	//2.解析日期范围参数
	var startDate, endDate time.Time
	layout := "2006-01-02"

	// 如果没有提供开始日期，默认为30天前
	if req.StartDate == "" {
		startDate = time.Now().AddDate(0, 0, -30)
	} else {
		startDate, err = time.Parse(layout, req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format, should be YYYY-MM-DD")
		}
	}

	// 如果没有提供结束日期，默认为今天
	if req.EndDate == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse(layout, req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format, should be YYYY-MM-DD")
		}
		// 确保结束日期包含整天的数据
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	//3.查询总点击量
	totalClicks, err := l.svcCtx.ClickStatisticsModel.CountTotalClicks(l.ctx, req.ShortUrl, startDate, endDate)
	if err != nil {
		logx.Errorw("CountTotalClicks failed", logx.LogField{Key: "err", Value: err.Error()})
		totalClicks = 0
	}

	//4.查询总唯一访问者数量
	totalUniqueVisitors, err := l.svcCtx.ClickStatisticsModel.CountTotalUniqueVisitors(l.ctx, req.ShortUrl, startDate, endDate)
	if err != nil {
		logx.Errorw("CountTotalUniqueVisitors failed", logx.LogField{Key: "err", Value: err.Error()})
		totalUniqueVisitors = 0
	}

	//5.查询每日统计数据
	dailyStatsRows, err := l.svcCtx.ClickStatisticsModel.GetDailyStats(l.ctx, req.ShortUrl, startDate, endDate)
	if err != nil {
		logx.Errorw("GetDailyStats failed", logx.LogField{Key: "err", Value: err.Error()})
		dailyStatsRows = []model.DailyStatsRow{}
	}

	//6.转换为响应格式
	dailyStats := make([]types.DailyStats, len(dailyStatsRows))
	for i, row := range dailyStatsRows {
		dailyStats[i] = types.DailyStats{
			Date:           row.ClickDate,
			Clicks:         row.TotalClicks,
			UniqueVisitors: row.UniqueVisitors,
		}
	}

	//7.构建并返回响应
	resp = &types.ShortUrlStatsResponse{
		TotalClicks:         totalClicks,
		TotalUniqueVisitors: totalUniqueVisitors,
		DailyStats:          dailyStats,
	}

	return resp, nil
}
