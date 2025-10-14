package main

import (
	"flag"
	"fmt"
	"shortener/pkg/base62"

	"shortener/internal/config"
	"shortener/internal/handler"
	"shortener/internal/svc"
	"shortener/internal/worker"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortener-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	//base62模块的初始化
	base62.MustInit(c.BaseString)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 初始化并启动点击统计worker
	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: c.CatheRedis[0].Host,
	})
	clickWorker := worker.NewClickWorker(redisClient, ctx.ClickStatisticsModel)
	clickWorker.Start()
	defer clickWorker.Stop()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
