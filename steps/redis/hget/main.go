package main

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/redis/base"
)

type Args struct {
	base.Args
	Key   string `env:"KEY,required"`
	Field string `env:"FIELD"`
}

type RedisHget struct {
	args        Args
	redisClient *redis.Client
}

func (l *RedisHget) Init() error {
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

func (l *RedisHget) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	cmd := l.redisClient.HGet(context.Background(), l.args.Key, l.args.Field)
	if cmd == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("invalid result returned from hget operation")
	}

	val, err := cmd.Result()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("failed getting field value: validate key is valid: %w", err)
	}

	gc.Set(val, "output")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&RedisHget{})
}
