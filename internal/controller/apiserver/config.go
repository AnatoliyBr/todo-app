package apiserver

import "os"

type Config struct {
	BindAddr  string `toml:"bind_addr"`
	SecretKey string
}

func NewConfig() *Config {
	sk := os.Getenv("SECRET_KEY")
	if sk == "" {
		sk = "Hj/XWYwIXa9JIXQmZzIHjoMuD/LQiA+omQjUhQx3QmHA/VcyJNFCn8btu3Vh3kGw"
	}

	return &Config{
		BindAddr:  ":8080",
		SecretKey: sk,
	}
}
