package cache

import (
	"auth_service/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthCache struct {
	redis *redis.Client
}

func NewAuthCache(redis *redis.Client) *AuthCache {
	return &AuthCache{redis: redis}
}

// CreateOrUpdateUserByEmail adds or updates a user in Redis by email.
func (a *AuthCache) CreateOrUpdateUserByEmail(ctx context.Context, user *models.User) error {
	key := fmt.Sprintf("email:%s", user.Email)
	data, err := json.Marshal(user)
	if err != nil {
		log.Println("Failed to marshal user: ", err)
		return err
	}

	_, err = a.redis.Set(ctx, key, data, 24*time.Hour).Result()
	if err != nil {
		log.Println("Failed to set user in Redis by email: ", err)
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user from Redis by email.
func (a *AuthCache) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	key := fmt.Sprintf("email:%s", email)
	cmd := a.redis.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Println("Redis connection error: ", err)
		return nil, err
	}

	user := models.User{}
	data, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&user); err != nil {
		log.Println("Failed decoding data: ", err)
		return nil, err
	}

	return &user, nil
}

// DeleteUserByEmail removes a user from Redis by email.
func (a *AuthCache) DeleteUserByEmail(ctx context.Context, email string) error {
	key := fmt.Sprintf("email:%s", email)
	_, err := a.redis.Del(ctx, key).Result()
	if err != nil {
		log.Println("Failed to delete user from Redis by email: ", err)
		return err
	}

	return nil
}

// CreateOrUpdateUserByID adds or updates a user in Redis by user ID.
func (a *AuthCache) CreateOrUpdateUserByID(ctx context.Context, user *models.User) error {
	key := fmt.Sprintf("user_id:%s", user.UserId)
	data, err := json.Marshal(user)
	if err != nil {
		log.Println("Failed to marshal user: ", err)
		return err
	}

	_, err = a.redis.Set(ctx, key, data, 24*time.Hour).Result()
	if err != nil {
		log.Println("Failed to set user in Redis by user ID: ", err)
		return err
	}

	return nil
}

// GetUserByID retrieves a user from Redis by user ID.
func (a *AuthCache) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	key := fmt.Sprintf("user_id:%s", userId)
	cmd := a.redis.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Println("Redis connection error: ", err)
		return nil, err
	}

	user := models.User{}
	data, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&user); err != nil {
		log.Println("Failed decoding data: ", err)
		return nil, err
	}

	return &user, nil
}

// DeleteUserByID removes a user from Redis by user ID.
func (a *AuthCache) DeleteUserByID(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user_id:%s", userId)
	_, err := a.redis.Del(ctx, key).Result()
	if err != nil {
		log.Println("Failed to delete user from Redis by user ID: ", err)
		return err
	}

	return nil
}
