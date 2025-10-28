package connection

import (
	"context"
	"fmt"
	"sanjay/config"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func RedisEnvCheck(r *config.RedisConfig) bool {
	if r.Addr == "" || r.Addr == "" || r.Password == "" || r.RedisHost == "" {
		return false
	}
	return true
}
func CheckRedisConnection(conn *redis.Client) error {
	_, err := conn.Ping(context.Background()).Result()
	return err
}

func NewRedisConnection() (*redis.Client, error) {
	r := config.LoadEnv().GetRedis()
	if exits := RedisEnvCheck(r); !exits {
		return nil, fmt.Errorf("Redis Env  variable missing")
	}
	db, _ := strconv.Atoi(r.DB)
	clint := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.RedisHost, r.Addr),
		Password: r.Password,
		DB:       db,
	})
	err := CheckRedisConnection(clint)
	return clint, err
}
