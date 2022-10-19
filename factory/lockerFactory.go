package factory

import (
	"github.com/ME-MotherEarth/me-notifier/common"
	"github.com/ME-MotherEarth/me-notifier/config"
	"github.com/ME-MotherEarth/me-notifier/disabled"
	"github.com/ME-MotherEarth/me-notifier/redis"
)

// CreateLockService creates lock service component based on config
func CreateLockService(checkDuplicates bool, config config.RedisConfig) (redis.LockService, error) {
	if !checkDuplicates {
		return disabled.NewDisabledRedlockWrapper(), nil
	}

	redisClient, err := createRedisClient(config)
	if err != nil {
		return nil, err
	}

	redlockArgs := redis.ArgsRedlockWrapper{
		Client:       redisClient,
		TTLInMinutes: config.TTL,
	}
	lockService, err := redis.NewRedlockWrapper(redlockArgs)
	if err != nil {
		return nil, err
	}

	return lockService, nil
}

func createRedisClient(cfg config.RedisConfig) (redis.RedLockClient, error) {
	switch cfg.ConnectionType {
	case common.RedisInstanceConnType:
		return redis.CreateSimpleClient(cfg)
	case common.RedisSentinelConnType:
		return redis.CreateFailoverClient(cfg)
	default:
		return nil, common.ErrInvalidRedisConnType
	}
}
