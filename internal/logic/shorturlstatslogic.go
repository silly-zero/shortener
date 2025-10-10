package logic

import (
	"context"

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

	return
}
