package connection

import (
	"context"
	"fmt"
	"sanjay/config"

	"github.com/redis/go-redis/v9"
)

func RedisEnvCheck(r *config.RedisConfig) bool {
	if r.ConnStr == " " {
		return false
	}
	return true
}
func CheckRedisConnection(conn *redis.Client) error {
	_, err := conn.Ping(context.Background()).Result()
	return err
}

func NewRedisConnection() (*redis.Client, error) {
	addr := config.LoadEnv().GetRedis()
	if exits := RedisEnvCheck(addr); !exits {
		return nil, fmt.Errorf("Redis Env  variable missing")
	}
	opt, _ := redis.ParseURL(addr.ConnStr)
	clint := redis.NewClient(opt)
	err := CheckRedisConnection(clint)
	return clint, err
}
