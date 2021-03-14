package base

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/sirupsen/logrus"
)

type Args struct {
	Token       string `env:"TOKEN,required"`
	UseEUApiURL bool   `env:"USE_EU_API" envDefault:"false"`
}

func Config(args Args) *client.Config {
	cfg := &client.Config{
		ApiKey:   args.Token,
		LogLevel: logrus.ErrorLevel,
	}

	if args.UseEUApiURL {
		cfg.OpsGenieAPIURL = client.API_URL_EU
	}

	return cfg
}
