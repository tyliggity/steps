package base

type Args struct {
	Host       string `env:"HOST,required"`
	Port       string `env:"PORT" envDefault:"8086"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
	Database   string `env:"DATABASE,required"`
	SSL        bool   `env:"SSL" envDefault:"false"`
	UnsafeSSL  bool   `env:"UNSAFE_SSL" envDefault:"false"`
	BinaryName string `env:"BINARY_NAME" envDefault:"influx"`
}
