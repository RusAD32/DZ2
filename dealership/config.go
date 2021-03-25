package dealership

type Config struct {
	//DatabaseURL ...
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
