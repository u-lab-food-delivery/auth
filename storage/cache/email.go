package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type EmailCache struct {
	redis *redis.Client
}

func NewEmailCache(redis *redis.Client) *EmailCache {
	return &EmailCache{redis: redis}
}

func (e *EmailCache) SaveLink(email, link string) error {
	key := fmt.Sprintf("link:%s", link)
	err := e.redis.Set(context.TODO(), key, email, time.Minute*2).Err()
	if err != nil {
		return err
	}

	return nil
}

func (e *EmailCache) GetEmailByLink(link string) (string, error) {
	key := fmt.Sprintf("link:%s", link)
	cmd := e.redis.Get(context.TODO(), key)
	if cmd.Err() == redis.Nil {
		return "", nil
	}
	email := cmd.Val()

	return email, nil
}
