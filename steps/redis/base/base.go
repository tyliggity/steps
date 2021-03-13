package base

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Args struct {
	RedisIP       string `env:"REDIS_IP"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisUsername string `env:"REDIS_USERNAME"`
	RedisUrl      string `env:"REDIS_URL"`
	Database      int    `env:"DATABASE" envDefault:"0"`
}

func InitRedisClient(args Args) (*redis.Client, error) {
	if args.RedisUrl == "" && args.RedisIP == "" {
		return nil, fmt.Errorf("invalid arguments, either REDIS_URL or REDIS_IP should be provided")
	}

	var err error
	var opts *redis.Options

	if args.RedisUrl != "" {
		opts, err = redis.ParseURL(args.RedisUrl)
		if err != nil {
			return nil, err
		}
	} else {
		opts = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", args.RedisIP, args.RedisPort),
			Password: args.RedisPassword,
			Username: args.RedisUsername,
			DB:       args.Database,
		}
	}

	return redis.NewClient(opts), nil
}

func ParseRESP(stepOutput []byte) ([]byte, error) {
	// RESP = REdis Serialization Protocol
	str := strings.TrimSuffix(string(stepOutput), "\n")
	rows := strings.Split(str, "\n")
	var m []map[string]string

	for _, row := range rows {
		entries := strings.Fields(row)
		mRow := make(map[string]string)

		for _, e := range entries {
			parts := strings.Split(e, "=")
			mRow[parts[0]] = parts[1]
		}

		m = append(m, mRow)
	}

	return json.Marshal(m)
}
