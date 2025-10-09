package sequence

import (
	"sync"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 批量Redis发号器
type BatchRedis struct {
	client     *redis.Redis
	key        string
	bathchSize int64
	currentMax int64
	currentVal int64
	lock       sync.Mutex
}

func NewBatchRedis(redisAddr string, batchSize int64) *BatchRedis {
	//初始化Redis客户端
	r := redis.MustNewRedis(redis.RedisConf{
		Host: redisAddr,
	})
	return &BatchRedis{
		client:     r,
		key:        "shortener:sequence",
		bathchSize: batchSize,
	}
}

func (r *BatchRedis) Next() (seq uint64, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return
}
