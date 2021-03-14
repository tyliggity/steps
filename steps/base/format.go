package main

import (
	"flag"
	"fmt"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/formats/json"
	"github.com/stackpulse/public-steps/common/formats/raw"
	"github.com/stackpulse/public-steps/common/gcs"
	"io/ioutil"
	"os"
	"strings"
)

type formatArgs struct {
	Input  string
	Format string
}

var gFormatArgs = &formatArgs{}

func defFormatArgs(flagSet *flag.FlagSet) {
	flagSet.StringVar(&gFormatArgs.Input, "input", "-", "input of json data (default is '-' means stdin)")
	flagSet.StringVar(&gFormatArgs.Format, "format", "", "the format of the output text")
}

func downloadGcsFile(url string) ([]byte, error) {
	bucket, object, err := gcs.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("gcs url parse failed: %w", err)
	}
	data, err := gcs.Download(bucket, object)
	if err != nil {
		return nil, fmt.Errorf("gcs download failed: %w", err)
	}
	return data, nil
}

func getData(input string) ([]byte, error) {
	if strings.HasPrefix(input, "gs://") {
		return downloadGcsFile(input)
	}

	inputFile := os.Stdin
	if input != "-" {
		var err error
		inputFile, err = os.Open(input)
		if err != nil {
			return nil, err
		}
	}

	return ioutil.ReadAll(inputFile)
}

func runFormat() error {
	data, err := getData(gFormatArgs.Input)
	if err != nil {
		return err
	}

	// First, printing the raw data
	fmt.Println(string(data))

	format := gFormatArgs.Format
	if format == "" {
		format = env.Formatter()
	}

	var output string

	switch format {
	case "json":
		output, err = json.Format(data)
	case "raw":
		output, err = raw.Format(data)
	case "print":
		output = "" // Already printed
	default:
		err = fmt.Errorf("unknown format %s", gFormatArgs.Format)
	}

	if err != nil {
		return fmt.Errorf("formatting error: %w", err)
	}

	fmt.Print(output)
	return nil
}
