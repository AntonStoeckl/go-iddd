package cmd

import (
	"os"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type Config struct {
	Postgres struct {
		DSN                    string
		MigrationsPathCustomer string
	}
	GRPC struct {
		HostAndPort string
	}
	REST struct {
		HostAndPort string
	}
}

func MustBuildConfigFromEnv(logger *Logger) *Config {
	var err error
	conf := &Config{}
	msg := "mustBuildConfigFromEnv: %s - Hasta la vista, baby!"

	if conf.Postgres.DSN, err = conf.stringFromEnv("POSTGRES_DSN"); err != nil {
		logger.Fatalf(msg, err)
	}

	if conf.Postgres.MigrationsPathCustomer, err = conf.stringFromEnv("POSTGRES_MIGRATIONS_PATH_CUSTOMER"); err != nil {
		logger.Fatalf(msg, err)
	}

	if conf.GRPC.HostAndPort, err = conf.stringFromEnv("GRPC_HOST_AND_PORT"); err != nil {
		logger.Fatalf(msg, err)
	}

	if conf.REST.HostAndPort, err = conf.stringFromEnv("REST_HOST_AND_PORT"); err != nil {
		logger.Fatalf(msg, err)
	}

	return conf
}

func (conf Config) stringFromEnv(envKey string) (string, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return "", errors.Mark(errors.Newf("config value [%s] missing in env", envKey), lib.ErrTechnical)
	}

	return envVal, nil
}
