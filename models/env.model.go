package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Config struct {
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	PGHOST       string `mapstructure:"PGHOST"`
	PGUSER       string `mapstructure:"PGUSER"`
	PGPASSWORD   string `mapstructure:"PGPASSWORD"`
	PGDATABASE   string `mapstructure:"PGDATABASE"`
	PGPORT       string `mapstructure:"PGPORT"`
}

func (config *Config) Validate() error {
	return validation.ValidateStruct(config,
		validation.Field(&config.DATABASE_URL, validation.Required, is.URL),
		validation.Field(&config.PGHOST, validation.Required),
		validation.Field(&config.PGUSER, validation.Required),
		validation.Field(&config.PGPASSWORD, validation.Required),
		validation.Field(&config.PGDATABASE, validation.Required),
		validation.Field(&config.PGPORT, validation.Required, is.Port),
	)
}
