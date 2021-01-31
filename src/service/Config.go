package service

import (
	"os"

	"github.com/AntonStoeckl/go-iddd/src/shared"
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

// This is also used by Config_test.go to check that all keys exist in Env,
// so always add new keys here!
var ConfigExpectedEnvKeys = map[string]string{
	"pgDSN":  "POSTGRES_DSN",
	"pgMPC":  "POSTGRES_MIGRATIONS_PATH_CUSTOMER",
	"grpcHP": "GRPC_HOST_AND_PORT",
	"restHP": "REST_HOST_AND_PORT",
}

func MustBuildConfigFromEnv(logger *shared.Logger) *Config {
	var err error
	conf := &Config{}
	msg := "mustBuildConfigFromEnv: %s - Hasta la vista, baby!"

	if conf.Postgres.DSN, err = conf.stringFromEnv(ConfigExpectedEnvKeys["pgDSN"]); err != nil {
		logger.Panicf(msg, err)
	}

	if conf.Postgres.MigrationsPathCustomer, err = conf.stringFromEnv(ConfigExpectedEnvKeys["pgMPC"]); err != nil {
		logger.Panicf(msg, err)
	}

	if conf.GRPC.HostAndPort, err = conf.stringFromEnv(ConfigExpectedEnvKeys["grpcHP"]); err != nil {
		logger.Panicf(msg, err)
	}

	if conf.REST.HostAndPort, err = conf.stringFromEnv(ConfigExpectedEnvKeys["restHP"]); err != nil {
		logger.Panicf(msg, err)
	}

	return conf
}

func (conf Config) stringFromEnv(envKey string) (string, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return "", errors.Mark(errors.Newf("config value [%s] missing in env", envKey), shared.ErrTechnical)
	}

	return envVal, nil
}
