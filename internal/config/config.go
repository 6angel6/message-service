package config

import "os"

type Config struct {
	HTTPAddr string
}

func Read() Config {
	var config Config
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		config.HTTPAddr = httpAddr
	}
	return config
}
