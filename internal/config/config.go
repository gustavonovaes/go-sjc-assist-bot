package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TELEGRAM_TOKEN       string `required:"true"`
	TELEGRAM_WEBHOOK_URL string `required:"true"`

	MONGODB_URI string `required:"true"`
}

func New() Config {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	return config
}
