package redis

import (
	"LipidClinic/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb *redis.Client
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.redis.New"

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(cfg.Redis.Host, ":", cfg.Redis.Port),
		Password: "",
		DB:       0,
	})

	return &Storage{rdb: rdb}, nil
}
