package logic

import (
	"context"
	"database/sql"
	"errors"

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
	// todo: add your logic here and delete this line
	//1.验证短链接是否存在
	if _, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true}); err != nil {
		return nil, errors.New("short url not found")
	}
	//2.查询短链接的统计信息

	return
}
