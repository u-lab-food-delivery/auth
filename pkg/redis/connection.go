package redis

import (
	"auth_service/config"

	"github.com/go-redis/redis/v8"
)

func ConnectDB(cnf *config.RedisConfig) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     cnf.Host + ":" + cnf.Port,
			Password: "",
			DB:       0,
		},
	)
}
