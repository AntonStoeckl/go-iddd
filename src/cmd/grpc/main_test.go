package main

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/cmd"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	customergrpc "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStartGRPCServer(t *testing.T) {
	logger := shared.NewNilLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	postgresDBConn := cmd.MustInitPostgresDB(config, logger)
	diContainer := cmd.MustBuildDIContainer(
		config,
		logger,
		cmd.UsePostgresDBConn(postgresDBConn),
		cmd.ReplaceGRPCCustomerServer(grpcCustomerServerStub()),
	)

	exitWasCalled := false
	exitFn := func() {
		exitWasCalled = true
	}

	terminateDelay := time.Millisecond * 100

	service := InitService(config, logger, exitFn, diContainer)

	Convey("Start the gRPC server as a goroutine", t, func() {
		go service.StartGRPCServer()

		Convey("gPRC server should handle requests", func() {
			client := customerGRPCClient(config)
			res, err := client.Register(context.Background(), &customergrpc.RegisterRequest{})
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.Id, ShouldNotBeEmpty)

			Convey(fmt.Sprintf("It should wait for stop signal (scheduled after %s)", terminateDelay), func() {
				start := time.Now()
				go func() {
					time.Sleep(terminateDelay)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				service.WaitForStopSignal()

				So(time.Now(), ShouldHappenOnOrAfter, start.Add(terminateDelay))

				Convey("Stop signal should issue Shutdown", func() {
					Convey("Shutdown should stop gRPC server", func() {
						_, err = client.Register(context.Background(), &customergrpc.RegisterRequest{})
						So(err, ShouldBeError)
						So(status.Code(err), ShouldResemble, codes.Unavailable)

						Convey("Shutdown should close PostgreSQL connection", func() {
							err := postgresDBConn.Ping()
							So(err, ShouldBeError)
							So(err.Error(), ShouldContainSubstring, "database is closed")

							Convey("Shutdown should call exit", func() {
								So(exitWasCalled, ShouldBeTrue)
							})
						})
					})
				})
			})
		})
	})
}

/*** Helper functions ***/

func grpcCustomerServerStub() customergrpc.CustomerServer {
	customerServer := customergrpc.NewCustomerServer(
		func(emailAddress, givenName, familyName string) (value.CustomerID, error) {
			return value.GenerateCustomerID(), nil
		},
		func(customerID, confirmationHash string) error {
			return nil
		},
		func(customerID, emailAddress string) error {
			return nil
		},
		func(customerID, givenName, familyName string) error {
			return nil
		},
		func(customerID string) error {
			return nil
		},
		func(customerID string) (customer.View, error) {
			return customer.View{}, nil
		},
	)

	return customerServer
}

func customerGRPCClient(config *cmd.Config) customergrpc.CustomerClient {
	grpcClientConn, _ := grpc.DialContext(context.Background(), config.GRPC.HostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	client := customergrpc.NewCustomerClient(grpcClientConn)

	return client
}
