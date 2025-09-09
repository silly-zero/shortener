package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB ShortURLDB

	Sequence struct {
		DSN string
	}
}
type ShortURLDB struct {
	DSN string
}
