package archivedata

type RedisOption struct {
	RedisAddr     string `flag:"redisaddr"`
	IDSeverAddr   string `flag:"idsvraddr"`
}

func NewRedisOption() *RedisOption {
	p := &RedisOption{
		RedisAddr:     "redis//:6379/4?pwd=probeKing+5688",
	}
	return p
}
