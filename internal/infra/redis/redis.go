package redis

import (
	goredis "github.com/redis/go-redis/v9"

	"friend_zone/internal/config"
)

func New(cfg config.RedisConfig) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
