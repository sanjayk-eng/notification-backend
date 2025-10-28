package service

import (
	"context"
	"encoding/json"
	"sanjay/api/util"
	"time"

	"github.com/redis/go-redis/v9"
)

const KEY = "otp"

type RedisCline struct {
	Redis *redis.Client
}

func NewRadishImplement(redis *redis.Client) *RedisCline {
	return &RedisCline{
		Redis: redis,
	}
}

func (r *RedisCline) StoreOTP(c context.Context, phoneNum, otp string) error {
	key := KEY + phoneNum
	return r.Redis.Set(c, key, otp, 1*time.Minute).Err()
}
func (r *RedisCline) GetOTP(c context.Context, phoneNum string) (string, error) {
	key := KEY + phoneNum
	return r.Redis.Get(c, key).Result()
}
func (r *RedisCline) DeleteOTP(c context.Context, phoneNum string) error {
	key := KEY + phoneNum
	return r.Redis.Del(c, key).Err()
}

// Store message in Redis list
func (r *RedisCline) PushMessage(ctx context.Context, msg util.WSMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	key := "messages:global"
	if err := r.Redis.LPush(ctx, key, data).Err(); err != nil {
		return err
	}

	// Keep only last 100 messages
	return r.Redis.LTrim(ctx, key, 0, 199).Err()
}

// Retrieve last 100 messages
func (r *RedisCline) GetMessages(ctx context.Context) ([]util.WSMessage, error) {
	key := "messages:global"
	list, err := r.Redis.LRange(ctx, key, 0, 199).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]util.WSMessage, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		var msg util.WSMessage
		if err := json.Unmarshal([]byte(list[i]), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

// Publish message to all subscribers
func (r *RedisCline) PublishMessage(ctx context.Context, msg util.WSMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.Redis.Publish(ctx, "chat:global", data).Err()
}

// Subscribe to global chat channel
func (r *RedisCline) SubscribeChannel(ctx context.Context) *redis.PubSub {
	return r.Redis.Subscribe(ctx, "chat:global")
}
