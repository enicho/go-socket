package util

import (
	"github.com/gomodule/redigo/redis"
)

var redisCli redis.Conn

func InitUtil() {
	redisCli, _ = redis.Dial("tcp", ":6379")
}

func GetClient() redis.Conn {
	if redisCli == nil {
		InitUtil()
	}

	return redisCli
}
