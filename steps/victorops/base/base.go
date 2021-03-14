package base

type Args struct {
	ApiID  string `env:"VICTOROPS_API_ID,required"`
	ApiKey string `env:"VICTOROPS_API_KEY,required"`
}
