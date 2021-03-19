package main

import (
	"fmt"
	"os"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/elastic/base"
)

func main() {
	cfg := &base.Args{}
	envconf.Parse(cfg)

	client, _ := base.CreateClient(cfg)
	info, err := client.Info()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	defer info.Body.Close()

	if info.IsError() {
		fmt.Printf("Error: %s", info.String())
	}

	fmt.Printf("%s", info)
	os.Exit(0)
}
