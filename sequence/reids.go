package sequence

// 基于Redis的取号器
type Reids struct {
}

func NewRedis(redisAddr string) *Reids {
	return &Reids{}
}

func (r *Reids) Next() (seq uint64, err error) {
	//使用redis实现发号器
	return
}
