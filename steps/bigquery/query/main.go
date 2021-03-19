package main

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
	"google.golang.org/api/iterator"
)

type Args struct {
	ProjectID string `env:"PROJECT_ID,required"`
	Query     string `env:"QUERY,required"`
}

type BQ struct {
	client *bigquery.Client
	args   Args
}

func (c *BQ) Init() error {
	err := envconf.Parse(&c.args)
	if err != nil {
		return err
	}
	c.client, err = bigquery.NewClient(context.Background(), c.args.ProjectID)
	return err
}

func (c *BQ) Run() (int, []byte, error) {
	defer c.client.Close()
	q := c.client.Query(c.args.Query)
	it, err := q.Read(context.Background())
	if err != nil {
		return 1, nil, err
	}

	var rows []map[string]bigquery.Value
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Log("Bigquery read: %v", err)
		}
		rows = append(rows, row)
	}
	out, err := json.Marshal(rows)
	return 0, out, err
}

func main() {
	step.Run(&BQ{})
}
