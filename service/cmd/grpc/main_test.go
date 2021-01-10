package main

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

func TestGRPCServer(t *testing.T) {
	logger := shared.NewNilLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	grpcCustomerServer := buildCustomerGRPCServer()
	diContainer := mustBuildDIContainer(grpcCustomerServer)
	noopExit := func() {}
	go startGRPCServer(config, logger, diContainer.GetGRPCServer(), noopExit)
	client := buildCustomerGRPCClient(config)

	Convey("It should handle a gRPC request", t, func() {
		res, err := client.Register(context.Background(), &customergrpc.RegisterRequest{})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Id, ShouldNotBeEmpty)
	})

	shutdown(logger, diContainer.GetGRPCServer(), diContainer.GetPostgresDBConn(), noopExit)
}

func TestWaitForStopSignal(t *testing.T) {
	logger := shared.NewNilLogger()

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	delay := time.Millisecond * 30
	start := time.Now()
	go func() {
		time.Sleep(delay)
		stopSignalChannel <- os.Interrupt
	}()

	shutdown := func() {}
	waitForStopSignal(stopSignalChannel, logger, shutdown)
	end := time.Now()

	Convey("It should wait for stop signal", t, func() {
		So(end, ShouldNotHappenWithin, delay, start)
	})
}

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

func mustBuildDIContainer(grpcCustomerServer customergrpc.CustomerServer) *cmd.DIContainer {
	diContainer, err := cmd.NewDIContainer(
		serialization.MarshalCustomerEvent,
		serialization.UnmarshalCustomerEvent,
		customer.BuildUniqueEmailAddressAssertions,
		cmd.WithGRPCCustomerServer(grpcCustomerServer),
	)
	if err != nil {
		panic(err)
	}

	return diContainer
}
