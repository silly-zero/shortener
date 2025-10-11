package handler

import (
	"context"
	"net/http"

	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端真实IP
		getClientIP := func(r *http.Request) string {
			// 检查是否有代理IP
			xForwardedFor := r.Header.Get("X-Forwarded-For")
			if xForwardedFor != "" {
				return xForwardedFor
			}
			// 检查RemoteAddr
			return r.RemoteAddr
		}
		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//参数规则校验
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorw("validator check failed", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取客户端信息
		clientIP := getClientIP(r)
		userAgent := r.Header.Get("User-Agent")
		referer := r.Header.Get("Referer")

		// 将客户端信息添加到context中
		ctx := context.WithValue(r.Context(), "ip", clientIP)
		ctx = context.WithValue(ctx, "userAgent", userAgent)
		ctx = context.WithValue(ctx, "referer", referer)

		l := logic.NewShowLogic(ctx, svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			//httpx.OkJsonCtx(r.Context(), w, resp)
			//返回重定向
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)
		}
	}
}