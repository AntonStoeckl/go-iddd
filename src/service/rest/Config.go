package rest

import (
	"os"
	"strconv"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

type Config struct {
	REST struct {
		HostAndPort             string
		GRPCDialHostAndPort     string
		GRPCDialTimeout         int
		SwaggerFilePathCustomer string
	}
}

// ConfigExpectedEnvKeys - This is also used by Config_test.go to check that all keys exist in Env,
// so always add new keys here!
var ConfigExpectedEnvKeys = map[string]string{
	"restHostAndPort":         "REST_HOST_AND_PORT",
	"grpcDialHostAndPort":     "GRPC_HOST_AND_PORT",
	"restGrpcDialTimeout":     "REST_GRPC_DIAL_TIMEOUT",
	"swiggerFilePathCustomer": "SWAGGER_FILE_PATH_CUSTOMER",
}

func MustBuildConfigFromEnv(logger *shared.Logger) *Config {
	var err error
	conf := &Config{}
	msg := "mustBuildConfigFromEnv: %s - Hasta la vista, baby!"

	if conf.REST.HostAndPort, err = conf.stringFromEnv(ConfigExpectedEnvKeys["restHostAndPort"]); err != nil {
		logger.Panic().Msgf(msg, err)
	}

	if conf.REST.GRPCDialHostAndPort, err = conf.stringFromEnv(ConfigExpectedEnvKeys["grpcDialHostAndPort"]); err != nil {
		logger.Panic().Msgf(msg, err)
	}

	if conf.REST.GRPCDialTimeout, err = conf.intFromEnv(ConfigExpectedEnvKeys["restGrpcDialTimeout"]); err != nil {
		logger.Panic().Msgf(msg, err)
	}

	if conf.REST.SwaggerFilePathCustomer, err = conf.stringFromEnv(ConfigExpectedEnvKeys["swiggerFilePathCustomer"]); err != nil {
		logger.Panic().Msgf(msg, err)
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

func (conf Config) intFromEnv(envKey string) (int, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return 0, errors.Mark(errors.Newf("config value [%s] missing in env", envKey), shared.ErrTechnical)
	}

	intEnvVal, err := strconv.Atoi(envVal)
	if err != nil {
		return 0, errors.Mark(errors.Newf("config value [%s] is not convertable to integer", envKey), shared.ErrTechnical)
	}

	return intEnvVal, nil
}
