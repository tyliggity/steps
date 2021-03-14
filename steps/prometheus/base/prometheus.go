package base

import (
	"fmt"
	prometheus "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"os"
)

type Args struct {
	URL   string `env:"URL,required"`
	Debug bool   `env:"DEBUG" envDefault:"false"`
}

func ErrorAndExit(err string) {
	fmt.Fprintf(os.Stderr, "%s", err)
	os.Exit(1)
}

func Client(args *Args) (prometheus.Client, error) {
	client, err := prometheus.NewClient(prometheus.Config{Address: args.URL})
	if err != nil {
		return nil, fmt.Errorf("failed creatung prometheus client: %w", err)
	}
	return client, nil
}

func API(args *Args) (promv1.API, error) {
	client, err := Client(args)
	if err != nil {
		return nil, err
	}

	return promv1.NewAPI(client), nil
}
