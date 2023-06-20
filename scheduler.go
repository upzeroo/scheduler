package scheduler

import (
	"github.com/go-redis/redis"
	"github.com/hibiken/asynq"
)

const (
	QueueRealtime = "realtime"
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type (
	// Scheduler -- main scheduler struct definition
	Scheduler struct {
		Config Config
	}

	// Config -- general scheduler config
	Config struct {
		RedisURL   string
		WorkersNum int
	}
)

// New -- scheduler construct here
func New(config Config) *Scheduler {
	return &Scheduler{
		Config: config,
	}
}

func (s *Scheduler) GetRedisClientOpts() (asynq.RedisClientOpt, error) {
	var opts asynq.RedisClientOpt

	redisURLOpts, err := redis.ParseURL(s.Config.RedisURL)
	if err != nil {
		return opts, err
	}

	opts = asynq.RedisClientOpt{
		Network:  redisURLOpts.Network,
		Addr:     redisURLOpts.Addr,
		Password: redisURLOpts.Password,
		DB:       redisURLOpts.DB,
		PoolSize: redisURLOpts.PoolSize,
	}

	return opts, nil
}
