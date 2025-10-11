package logic

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/logx"
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
	//根据短链接查询长链接
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{Valid: true, String: req.ShortUrl})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("404")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Value: err.Error(), Key: "err"})
		return nil, err
	}

	// 保存点击记录（异步处理，不阻塞主流程）
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 获取请求中的IP和User-Agent
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

		// 保存点击记录
		_, err := l.svcCtx.ClickStatisticsModel.Insert(l.ctx, &model.ClickStatistics{
			Surl:      req.ShortUrl,
			ClickTime: time.Now(),
			Ip:        sql.NullString{String: ip, Valid: true},
			UserAgent: sql.NullString{String: userAgent, Valid: true},
			Referer:   sql.NullString{String: referer, Valid: referer != ""},
		})
		if err != nil {
			logx.Errorw("保存点击记录失败", logx.LogField{Key: "err", Value: err.Error()})
		}
	}()

	//返回长链接
	return &types.ShowResponse{
			LongUrl: u.Lurl.String,
		},
		nil
}
