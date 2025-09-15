package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	urltool "shortener/pkg/url"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转换长链接为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	//1.校验输入数据
	//1.1 数据不能为空
	//使用validator校验输入数据
	//1.2 长链接是个可以请求的网址
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("无效链接")
	}
	//1.3 判断长链接是否已经存在
	//1.3.1 给长链接生成md5
	md5Value := md5.Sum([]byte(req.LongUrl))
	//1.3.2 检查md5是否已经存在
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该链接已存在为%s", u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//1.4 输入的不能是个短链接
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "lUrl", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该链接已经是短链接")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	var short string
	for {
		//2.取号 基于mysql实现的发号器
		//每来一个转链请求，我们就从mysql使用replace into的发号表中取一个号
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		fmt.Println(seq)
		//3.号码转为短链接
		//3.1 安全性
		short = base62.IntToBase62(seq)
		//3.2 短链接怎么避免敏感的字符
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; ok {
			break // 避免敏感字符
		}
	}
	fmt.Printf("short: %s\n", short)
	//4.存储长链接和短链接的映射关系
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Surl: sql.NullString{String: short, Valid: true},
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
		},
	); err != nil {
		logx.Errorw("ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
	}
	//5.返回短链接
	//5.1 返回的是短域名+短链接 baidu.com/123
	shortUrl := l.svcCtx.Config.ShortDoamin + "/" + short
	return &types.ConvertResponse{
		ShortUrl: shortUrl,
	}, nil
}
