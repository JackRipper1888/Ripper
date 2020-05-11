package utils

import (
	"testing"
)

const (
	RedisStreamHM = "wdhqwuh"
)

func TestInitRedis(t *testing.T) {
	//初始化redis数据库
	InitRedis(
		RedisStreamHM,
		"183.60.143.82:36986",
		"YHuI=%(862hJKSGtuEq",
		0,
		500,
	)
	pipe := GetRedisPipe(RedisStreamHM).TxPipeline()
	//pipe.HSetNX("hx","117.176.227.111:23089","ashdiqwednjnb")
	pipe.HDel("hx", "117.176.227.111:23089")
	pipe.Exec()
}
