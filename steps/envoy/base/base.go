package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/common/step"
)

type Args struct {
	Host      string `env:"HOST,required"`
	Port      string `env:"PORT" envDefault:"8001"`
	SSL       bool   `env:"SSL" envDefault:"false"`
	UnsafeSSL bool   `env:"UNSAFE_SSL" envDefault:"false"`
	Endpoint  string `env:"ENDPOINT",required`
}

func (a *Args) GetURI() *url.URL {
	u := url.URL{
		Host:   net.JoinHostPort(a.Host, a.Port),
		Scheme: "http",
	}

	if a.SSL {
		u.Scheme = "https"
	}

	return &u
}

const (
	httpTimeout           = 10 * time.Second
	tlsHandshakeTimeout   = 10 * time.Second
	responseHeaderTimeout = 10 * time.Second
	expectContinueTimeout = 1 * time.Second
)

func newHttpClient(a *Args) *http.Client {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{},
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
		ExpectContinueTimeout: expectContinueTimeout,
	}

	if a.UnsafeSSL {
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

	return &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
	}
}

type EnvoyStep struct {
	args   Args
	client *http.Client
}

func (p *EnvoyStep) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	p.client = newHttpClient(&p.args)

	return nil
}

func (p *EnvoyStep) Run() (int, []byte, error) {

	targetUrl := fmt.Sprintf("%s%s", p.args.GetURI().String(), p.args.Endpoint)

	resp, err := p.client.Get(targetUrl)
	if err != nil {
		return 1, nil, fmt.Errorf("failed requesting url %s: %w", targetUrl, err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return 1, nil, fmt.Errorf("failed reading body: %w", err)
	}

	return 0, respBody, nil
}

func main() {
	step.Run(&EnvoyStep{})
}
