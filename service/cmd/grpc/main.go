package main

import (
	"database/sql"
	"net"
	"os"
	"os/signal"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

func main() {
	logger := shared.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)

	postgresDBConn := cmd.MustInitPostgresDB(config, logger)

	diContainer := cmd.MustBuildDIContainer(
		config,
		logger,
		serialization.MarshalCustomerEvent,
		serialization.UnmarshalCustomerEvent,
		customer.BuildUniqueEmailAddressAssertions,
		cmd.WithPostgresDBConn(postgresDBConn),
	)

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	shutdown := func() {
		shutdown(logger, diContainer.GetGRPCServer(), diContainer.GetPostgresDBConn(), osExit)
	}

	go startGRPCServer(config, logger, diContainer.GetGRPCServer(), shutdown)

	waitForStopSignal(stopSignalChannel, logger, shutdown)
}

func startGRPCServer(
	config *cmd.Config,
	logger *shared.Logger,
	grpcServer *grpc.Server,
	shutdown func(),
) {

	logger.Info("configuring gRPC server ...")

	listener, err := net.Listen("tcp", config.GRPC.HostAndPort)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		shutdown()
	}

	logger.Infof("starting gRPC server listening at %s ...", config.GRPC.HostAndPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		shutdown()
	}
}

func waitForStopSignal(
	stopSignalChannel chan os.Signal,
	logger *shared.Logger,
	shutdown func(),
) {

	logger.Info("start waiting for stop signal ...")

	sig := <-stopSignalChannel

	switch sig.(type) {
	case os.Signal:
		logger.Infof("received '%sig'", sig)
		close(stopSignalChannel)
		shutdown()
	}
}

func shutdown(
	logger *shared.Logger,
	grpcServer *grpc.Server,
	postgresDBConn *sql.DB,
	exit func(),
) {

	logger.Info("shutdown: stopping services ...")

	if grpcServer != nil {
		logger.Info("shutdown: stopping gRPC server gracefully ...")
		grpcServer.GracefulStop()
	}

	if postgresDBConn != nil {
		logger.Info("shutdown: closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("shutdown: failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	exit()
}

func osExit() {
	os.Exit(0)
}
