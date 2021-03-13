package main

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/redis/base"
)

type Args struct {
	base.Args
	Pattern string `env:"PATTERN,required"`
}

type RedisKeys struct {
	args        Args
	redisClient *redis.Client
}

func (l *RedisKeys) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	redisClient, err := base.InitRedisClient(l.args.Args)
	if err != nil {
		return err
	}

	l.redisClient = redisClient

	return nil
}

func (l *RedisKeys) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	cmd := l.redisClient.Keys(context.Background(), l.args.Pattern)
	if cmd == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("invalid result returned from keys operation")
	}

	val, err := cmd.Result()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("failed getting keys: validate pattern is valid: %w", err)
	}

	gc.Set(val, "output")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&RedisKeys{})
}
