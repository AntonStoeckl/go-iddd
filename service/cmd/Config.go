package cmd

import (
	"go-iddd/service/lib"
	"os"

	"github.com/cockroachdb/errors"
)

type Config struct {
	Postgres struct {
		DSN            string
		MigrationsPath string
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

	if conf.Postgres.MigrationsPath, err = conf.fromEnv("POSTGRES_MIGRATIONS_PATH"); err != nil {
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
