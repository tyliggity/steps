package main

import (
	"context"
	"encoding/json"
	"fmt"

	envconf "github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stackpulse/public-steps/redis/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	base.Args
	LastEntries int64 `env:"LAST_ENTRIES" envDefault:"10"`
}

type RedisSlowLog struct {
	args        Args
	redisClient *redis.Client
}

func (l *RedisSlowLog) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	l.redisClient, err = base.InitRedisClient(l.args.Args)

	return err
}

func (l *RedisSlowLog) Run() (exitCode int, output []byte, err error) {
	cmd := l.redisClient.SlowLogGet(context.Background(), l.args.LastEntries)
	if cmd == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("invalid result returned from slowlog operation")
	}

	val, err := json.Marshal(cmd.Val())
	if err != nil {
		return step.ExitCodeFailure, val, err
	}

	return step.ExitCodeOK, val, nil
}

func main() {
	step.Run(&RedisSlowLog{})
}
