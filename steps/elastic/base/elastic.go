package base

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	elastic "github.com/elastic/go-elasticsearch/v8"
)

const (
	CertFile    = "/tmp/cert-file"
	CertKeyFile = "/tmp/cert-file"
)

func createCertificate(certFile, certKey string) (*tls.Certificate, error) {
	// Write content to files
	err := ioutil.WriteFile(CertFile, []byte(certFile), 0644)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(CertKeyFile, []byte(certKey), 0644)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(CertFile, CertKeyFile)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

type Args struct {
	Host           string `env:"HOST,required"`
	Port           string `env:"PORT,required"`
	Username       string `env:"USERNAME"`
	Password       string `env:"PASSWORD"`
	CertContent    string `env:"CERT_CONTENT"`
	CertKeyContent string `env:"CERT_KEY_CONTENT"`
	UnsafeSSL      bool   `env:"UNSAFE_SSL"`
}

func CreateClient(args *Args) (*elastic.Client, error) {
	address := args.Host + ":" + args.Port
	cfg := &elastic.Config{
		Addresses: []string{
			address,
		},
		Username: args.Username,
		Password: args.Password,
	}

	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: args.UnsafeSSL,
		},
	}

	var cert *tls.Certificate
	var err error
	if !args.UnsafeSSL && args.CertContent != "" && args.CertKeyContent != "" {
		cert, err = createCertificate(args.CertContent, args.CertKeyContent)
		if err != nil {
			return nil, err
		}
		httpTransport.TLSClientConfig.Certificates = []tls.Certificate{*cert}
	}

	cfg.Transport = httpTransport

	client, err := elastic.NewClient(*cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}
