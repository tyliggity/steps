package base

type Args struct {
	Token    string `env:"TOKEN,required"`
	Endpoint string `env:"ENDPOINT" envDefault:"https://api.logz.io"`
}
