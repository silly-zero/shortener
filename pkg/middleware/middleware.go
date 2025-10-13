package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/chain"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/limit"
)

// WithRateLimiter 限流中间件
func WithRateLimiter(qps int) func(next http.Handler) http.Handler {
	l := limit.NewTokenLimiter(qps, qps*60, nil, "rate_limiter")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.Allow() {
				http.Error(w, "请求频率过快，请稍后重试", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// WithRequestContext 添加请求上下文信息
func WithRequestContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			//提取id
			ip := r.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}

			//处理代理
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				if len(ips) > 0 {
					ip = strings.TrimSpace(ips[0])
				}
			}
			//添加到上下文
			ctx = context.WithValue(ctx, "ip", ip)
			ctx = context.WithValue(ctx, "user_agent", r.UserAgent())
			ctx = context.WithValue(ctx, "referer", r.Referer())
			ctx = context.WithValue(ctx, "startTime", time.Now())

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}

}

// WithAccessLog 添加访问日志中间件
func WithAccessLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 调用下一个处理器
			next.ServeHTTP(w, r)

			// 记录访问日志
			duration := time.Since(start)
			logx.Infow("HTTP Request",
				logx.Field("method", r.Method),
				logx.Field("path", r.URL.Path),
				logx.Field("duration", duration),
				logx.Field("ip", r.RemoteAddr),
			)
		})
	}
}

// Chain 创建中间件链
func Chain(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return chain.New(handlers...)
}
