package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB ShortURLDB

	Sequence struct {
		DSN string
	}

	BaseString string

	ShortUrlBlackList []string
	ShortDoamin       string
	CatheRedis        cache.CacheConf //redis缓存
}
type ShortURLDB struct {
	DSN string
}
