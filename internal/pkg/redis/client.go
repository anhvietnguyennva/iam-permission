package redis

import (
	"context"
	"errors"
	"strings"

	e "github.com/anhvietnguyennva/go-error/pkg/errors"
	"github.com/go-redis/redis/v8"

	"iam-permission/internal/pkg/config"
)

var redisClient redis.UniversalClient

func InitClient() error {
	if redisClient == nil {
		cfg := config.Instance().Redis
		redisAddresses := strings.Split(cfg.Addresses, ",")
		if len(redisAddresses) == 0 {
			return errors.New("redis host is empty")
		}

		redisClient = newRedisClient(redisAddresses, cfg.MasterName, cfg.Password)
		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			return e.NewInfraErrorRedisConnect(err)
		}
	}
	return nil
}

func ClientInstance() redis.UniversalClient {
	return redisClient
}

func newRedisClient(redisAddresses []string, masterName string, pwd string) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      redisAddresses,
		MasterName: masterName,
		Password:   pwd,
	})
}
