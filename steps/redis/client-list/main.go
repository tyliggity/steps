package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	envconf "github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/redis/base"
)

type Args struct {
	base.Args
	OrderBy   string `env:"ORDER_BY"`
	OrderDesc bool   `env:"ORDER_DESC" envDefault:"false"`
	Limit     int    `env:"LIMIT" envDefault:"-1"`
}

type RedisClientList struct {
	args        Args
	redisClient *redis.Client
}

func (l *RedisClientList) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	if l.args.Limit < -1 {
		return fmt.Errorf("LIMIT cannot be a negative number (%v)", l.args.Limit)
	}

	l.redisClient, err = base.InitRedisClient(l.args.Args)

	return err
}

func (l *RedisClientList) order(list []map[string]string) []map[string]string {
	orderBy := l.args.OrderBy
	isDesc := l.args.OrderDesc

	sort.Slice(list, func(i, j int) bool {
		return (list[i][orderBy] < list[j][orderBy]) == isDesc
	})

	return list
}

func (l *RedisClientList) Run() (exitCode int, output []byte, err error) {
	cmd := l.redisClient.ClientList(context.Background())
	if cmd == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("invalid result returned from client list operation")
	}

	val, err := cmd.Bytes()
	if err != nil {
		return step.ExitCodeFailure, val, err
	}

	val, err = base.ParseRESP(val)
	if err != nil {
		return step.ExitCodeFailure, val, err
	}

	if l.args.OrderBy != "" || l.args.Limit != -1 {
		var list []map[string]string
		err := json.Unmarshal(val, &list)
		if err != nil {
			return step.ExitCodeFailure, val, err
		}

		if l.args.OrderBy != "" {
			list = l.order(list)
		}

		if l.args.Limit > -1 && l.args.Limit < len(list) {
			list = list[:l.args.Limit]
		}

		val, err = json.Marshal(&list)
		if err != nil {
			return step.ExitCodeFailure, val, err
		}
	}

	return step.ExitCodeOK, val, nil
}

func main() {
	step.Run(&RedisClientList{})
}
