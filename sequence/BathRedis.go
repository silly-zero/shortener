package sequence

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 批量Redis发号器
type BatchRedis struct {
	client     *redis.Redis
	key        string
	batchSize int64
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
		batchSize: batchSize,
	}
}

func (r *BatchRedis) Next() (seq uint64, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 如果当前批次号用完了，从Redis批量获取新的号段
	if r.currentVal >= r.currentMax {
		// 使用Redis的INCRBY命令原子性地获取一个批次的序列号
		newMax, err := r.client.Incrby(r.key, r.batchSize)
		if err != nil {
			logx.Errorw("client.IncrBy failed", logx.Field("err", err.Error()))
			return 0, err
		}

		// 更新批次范围
		r.currentMax = newMax
		r.currentVal = newMax - r.batchSize + 1
	}

	// 返回当前序列号并递增
	seq = uint64(r.currentVal)
	r.currentVal++

	return seq, nil
}