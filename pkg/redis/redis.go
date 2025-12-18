package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/supportbot/supportbot-go/internal/config"
)

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return client, nil
}

