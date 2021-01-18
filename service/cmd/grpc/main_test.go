package main

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/service/shared"
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
		cmd.ReplaceGRPCCustomerServer(buildCustomerGRPCServer()),
		cmd.UsePostgresDBConn(postgresDBConn),
	)
	grpcServer := diContainer.GetGRPCServer()

	exitWasCalled := false
	exit := func() {
		exitWasCalled = true
	}
	myShutdown := func() {
		shutdown(logger, grpcServer, postgresDBConn, exit)
	}

	terminateDelay := time.Millisecond * 100

	Convey("Start the gRPC server as a goroutine", t, func() {
		go startGRPCServer(config, logger, grpcServer, myShutdown)

		Convey(fmt.Sprintf("Schedule stop signal to be sent after %s", terminateDelay), func() {
			start := time.Now()
			go func() {
				time.Sleep(terminateDelay)
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			}()

			Convey("gPRC server should handle a request", func() {
				client := buildCustomerGRPCClient(config)
				res, err := client.Register(context.Background(), &customergrpc.RegisterRequest{})
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.Id, ShouldNotBeEmpty)

				Convey("Start waiting for stop signal", func() {
					waitForStopSignal(logger, myShutdown)

					Convey("It should wait for stop signal", func() {
						So(time.Now(), ShouldNotHappenWithin, terminateDelay, start)

						Convey("Once stop signal is received, it should call shutdown", func() {
							Convey("Shutdown should stop gRPC server", func() {
								_, err = client.Register(context.Background(), &customergrpc.RegisterRequest{})
								So(err, ShouldBeError)
								So(status.Code(err), ShouldResemble, codes.Unavailable)

								Convey("Shutdown should close PostgreSQL connection", func() {
									err := postgresDBConn.Ping()
									So(err, ShouldBeError)
									So(err.Error(), ShouldEqual, "sql: database is closed")

									Convey("Shutdown should call exit", func() {
										So(exitWasCalled, ShouldBeTrue)
									})
								})
							})
						})
					})
				})
			})
		})
	})
}

/*** Helper functions ***/

func buildCustomerGRPCServer() customergrpc.CustomerServer {
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

func buildCustomerGRPCClient(config *cmd.Config) customergrpc.CustomerClient {
	grpcClientConn, _ := grpc.DialContext(context.Background(), config.GRPC.HostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	client := customergrpc.NewCustomerClient(grpcClientConn)

	return client
}
