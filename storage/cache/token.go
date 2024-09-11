package cache

import (
	"auth_service/config"
	"auth_service/models"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type TokenCache struct {
	redis *redis.Client
}

func NewTokenCache(client *redis.Client) *TokenCache {
	return &TokenCache{
		redis: client,
	}
}

type RevokeTokens struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (t *TokenCache) RevokeToken(ctx context.Context, tokens *RevokeTokens) error {
	key := fmt.Sprintf("revoked_token:refresh_token:%s:access_token:%s", tokens.RefreshToken, tokens.AccessToken)

	claims, err := ExtractClaims(tokens.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to extract claims: %v", err)
	}

	expirationTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	timeLeft := expirationTime.Sub(time.Now())
	if timeLeft > 0 {
		err = t.redis.Set(ctx, key, "revoked", timeLeft).Err()
		if err != nil {
			return fmt.Errorf("failed to revoke token: %v", err)
		}
	}

	return nil
}

func (t *TokenCache) IsRevokedToken(ctx context.Context, refreshToken string, accessToken string) (bool, error) {
	key := fmt.Sprintf("revoked_token:refresh_token:%s:access_token:%s", refreshToken, accessToken)

	result, err := t.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check revoked token: %v", err)
	}

	return result == 1, nil
}
func ExtractClaims(tokenStr string) (*models.Claims, error) {
	claims := &models.Claims{} // Initialize the claims object
	cnf := config.NewConfig()
	cnf.Load()

	jwtKey := cnf.JWT.SecretKey

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to extract claims: %v", err)
	}

	return claims, nil
}
