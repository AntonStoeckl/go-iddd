package grpc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMustBuildConfigFromEnv(t *testing.T) {
	logger := shared.NewNilLogger()

	Convey("Given all values are set in Env", t, func() {
		Convey("When MustBuildConfigFromEnv is invoked", func() {
			config := grpc.MustBuildConfigFromEnv(logger)

			Convey("Then it should succeed", func() {
				wrapper := func() { grpc.MustBuildConfigFromEnv(logger) }
				So(wrapper, ShouldNotPanic)
				So(config, ShouldNotBeZeroValue)
			})
		})
	})

	for _, envKey := range grpc.ConfigExpectedEnvKeys {
		currentEnvKey := envKey

		Convey(fmt.Sprintf("Given %s is missing in Env", envKey), t, func() {
			origEnvVal := os.Getenv(currentEnvKey)
			err := os.Unsetenv(currentEnvKey)
			So(err, ShouldBeNil)

			Convey("When MustBuildConfigFromEnv is invoked", func() {
				wrapper := func() { grpc.MustBuildConfigFromEnv(logger) }

				Convey("It should panic", func() {
					So(wrapper, ShouldPanic)
				})
			})

			err = os.Setenv(currentEnvKey, origEnvVal)
			So(err, ShouldBeNil)
		})
	}
}
