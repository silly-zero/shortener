package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"
)

func ShortUrlStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortUrlStatsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShortUrlStatsLogic(r.Context(), svcCtx)
		resp, err := l.ShortUrlStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
