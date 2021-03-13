package base

import (
	"crypto/tls"
	"fmt"
	"github.com/Jeffail/gabs"
	envconf "github.com/caarlos0/env/v6"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Args struct {
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	Schema   string `env:"SCHEMA" envDefault:"https"`
	Verify   bool   `env:"VERIFY" envDefault:"false"`
}

func (a *Args) BaseArgs() *Args {
	return a
}

type BaseArgs interface {
	BaseArgs() *Args
}

type RabbitMQBase struct {
	*Args
}

func NewRabbitMQBase(args BaseArgs) (*RabbitMQBase, error) {
	if args == nil {
		args = &Args{}
	}
	if err := envconf.Parse(args); err != nil {
		return nil, err
	}

	return &RabbitMQBase{Args: args.BaseArgs()}, nil
}

func (r *RabbitMQBase) buildURL(api string) string {
	var url strings.Builder
	url.WriteString(r.Schema)
	url.WriteString("://")
	url.WriteString(r.Host)
	url.WriteString(":")
	url.WriteString(strconv.Itoa(r.Port))
	url.WriteString("/api/")
	url.WriteString(api)

	return url.String()
}

func (r *RabbitMQBase) RunQuery(api string) ([]byte, error) {
	client := &http.Client{}
	if !r.Verify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	url := r.buildURL(api)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(r.Username, r.Password)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending HTTP GET request to %s: %w", url, err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	parsed, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse output as JSON: %w", err)
	}

	retGc := gabs.New()
	_, _ = retGc.Set(parsed.Data(), "output")

	return retGc.Bytes(), nil
}
