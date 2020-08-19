package test

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/goapt/redis"
)

func NewRedis() *redis.Redis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return redis.NewRedis(redis.Open(s.Addr()))
}
