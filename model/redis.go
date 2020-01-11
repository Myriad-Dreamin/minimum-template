package model

import (
	"github.com/Myriad-Dreamin/minimum-template/config"
	"github.com/Myriad-Dreamin/minimum-lib/module"
	splayer "github.com/Myriad-Dreamin/minimum-template/model/sp-layer"
	"github.com/gomodule/redigo/redis"
	"time"
)

func OpenRedis(cfg *config.ServerConfig) (*redis.Pool, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				cfg.RedisConfig.ConnectionType, cfg.RedisConfig.Host,
				redis.DialPassword(cfg.RedisConfig.Password),
				redis.DialDatabase(cfg.RedisConfig.Database),
				redis.DialConnectTimeout(cfg.RedisConfig.ConnectionTimeout),
				redis.DialReadTimeout(cfg.RedisConfig.ReadTimeout),
				redis.DialWriteTimeout(cfg.RedisConfig.WriteTimeout),
				redis.DialKeepAlive(time.Minute*5),
			)
		},
		//TestOnBorrow:    nil,
		MaxIdle:     cfg.RedisConfig.MaxIdle,
		MaxActive:   cfg.RedisConfig.MaxActive,
		IdleTimeout: cfg.RedisConfig.IdleTimeout,
		Wait:        cfg.RedisConfig.Wait,
		//MaxConnLifetime: 0,
	}
	return pool, nil
}

func RegisterRedis(dep module.Module) bool {
	return splayer.RegisterRedis(dep)
}
