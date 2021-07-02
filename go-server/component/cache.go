package component

import (
	"fmt"

	"go-server/library/clean"
	"go-server/library/redis"
)

var CacheContainer *redis.ClientContainer

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
	PoolSize int    `env:"REDIS_POOLSIZE"`
}

func SetupCache() (err error) {
	CacheContainer, err = redis.NewContainer(getRedisConf)
	if err != nil {
		err = fmt.Errorf("redis.NewContainer: %w", err)
		return
	}
	clean.Push(CacheContainer)
	Conf.PushUpdater(CacheContainer)

	return
}

func getRedisConf() (cf *redis.ClientConf, err error) {
	cfg := &RedisConfig{}

	if err = Conf.Scan(cfg, "env"); err != nil {
		err = fmt.Errorf("Conf.Scan: %w", err)
		return
	}

	cf = &redis.ClientConf{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	}

	return
}
