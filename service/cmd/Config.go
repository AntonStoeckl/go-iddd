package cmd

import (
	"os"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type Config struct {
	Postgres struct {
		DSN                      string
		MigrationsPathEventstore string
		MigrationsPathCustomer   string
	}
}

func NewConfigFromEnv() (*Config, error) {
	var err error
	conf := &Config{}
	wrapMsg := "config.New"

	// Postgres
	if conf.Postgres.DSN, err = conf.fromEnv("POSTGRES_DSN"); err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	if conf.Postgres.MigrationsPathEventstore, err = conf.fromEnv("POSTGRES_MIGRATIONS_PATH_EVENTSTORE"); err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	if conf.Postgres.MigrationsPathCustomer, err = conf.fromEnv("POSTGRES_MIGRATIONS_PATH_CUSTOMER"); err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return conf, nil
}

func (conf Config) fromEnv(envKey string) (string, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return "", errors.Mark(errors.Newf("config value [%s] missing in env", envKey), lib.ErrTechnical)
	}

	return envVal, nil
}
