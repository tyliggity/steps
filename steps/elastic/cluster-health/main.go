package main

import (
	"bytes"
	"fmt"
	envconf "github.com/caarlos0/env/v6"
	elastic "github.com/stackpulse/public-steps/elastic/base"
	"os"
)

func main() {
	cfg := &elastic.Args{}
	envconf.Parse(cfg)

	client, _ := elastic.CreateClient(cfg)
	nodes, err := client.Cluster.Health()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(nodes.Body)

	fmt.Printf("%s", buf)
	os.Exit(0)
}
