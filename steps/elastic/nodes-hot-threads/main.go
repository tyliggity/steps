package main

import (
	"bytes"
	"fmt"
	"os"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/elastic/base"
)

func main() {
	cfg := &base.Args{}
	envconf.Parse(cfg)

	client, _ := base.CreateClient(cfg)
	nodesUsg, err := client.Nodes.HotThreads()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(nodesUsg.Body)

	fmt.Printf("%s", buf)
	os.Exit(0)
}
