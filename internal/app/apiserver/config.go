package apiserver

import "DZ2/dealership"

//General config for rest api
type Config struct {
	//Port for start api
	BindAddr   string `toml:"bind_addr"`
	LogLevel   string `toml:"log_level"`
	Dealership *dealership.Config
}

//Should return default config
func NewConfig() *Config {
	return &Config{
		BindAddr:   ":8080",
		LogLevel:   "debug",
		Dealership: dealership.NewConfig(),
	}
}
