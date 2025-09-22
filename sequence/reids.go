package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 基于Redis的取号器
type Redis struct {
	client *redis.Redis
	key    string
}

func NewRedis(redisAddr string) *Redis {
	//初始化Redis客户端
	r := redis.MustNewRedis(redis.RedisConf{
		Host: redisAddr,
	})
	return &Redis{
		client: r,
		key:    "shortener:sequence",
	}
}

func (r *Redis) Next() (seq uint64, err error) {
	//使用redis实现发号器
	//使用redis的incr实现原子递增
	lid, err := r.client.Incr(r.key)
	if err != nil {
		logx.Errorw("client.Incr failed", logx.Field("err", err.Error()))
		return 0, err
	}
	return uint64(lid), nil
}
