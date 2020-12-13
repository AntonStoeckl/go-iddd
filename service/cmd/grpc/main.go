package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var err error

	logger := shared.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)

	diContainer, err := cmd.Bootstrap(config, logger)
	if err != nil {
		shutdown(logger, diContainer, nil)
	}

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	grpcServer := buildGRPCServer(diContainer)

	go mustStartGRPC(config, logger, diContainer, grpcServer)

	waitForStopSignal(stopSignalChannel, logger, diContainer, grpcServer)
}

func buildGRPCServer(diContainer *cmd.DIContainer) *grpc.Server {
	grpcServer := grpc.NewServer()
	customergrpc.RegisterCustomerServer(grpcServer, diContainer.GetCustomerGRPCServer())
	reflection.Register(grpcServer)

	return grpcServer
}

func mustStartGRPC(
	config *cmd.Config,
	logger *shared.Logger,
	diContainer *cmd.DIContainer,
	grpcServer *grpc.Server,
) {

	logger.Info("configuring gRPC server ...")

	listener, err := net.Listen("tcp", config.GRPC.HostAndPort)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		shutdown(logger, diContainer, grpcServer)
	}

	logger.Infof("starting gRPC server listening at %s ...", config.GRPC.HostAndPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		shutdown(logger, diContainer, grpcServer)
	}
}

func waitForStopSignal(
	stopSignalChannel chan os.Signal,
	logger *shared.Logger,
	diContainer *cmd.DIContainer,
	grpcServer *grpc.Server,
) {

	logger.Info("start waiting for stop signal ...")

	s := <-stopSignalChannel

	switch s.(type) {
	case os.Signal:
		logger.Infof("received '%s'", s)
		close(stopSignalChannel)
		shutdown(logger, diContainer, grpcServer)
	}
}

func shutdown(
	logger *shared.Logger,
	diContainer *cmd.DIContainer,
	grpcServer *grpc.Server,
) {

	logger.Info("shutdown: stopping services ...")

	if grpcServer != nil {
		logger.Info("shutdown: stopping gRPC server gracefully ...")
		grpcServer.GracefulStop()
	}

	if postgresDBConn := diContainer.GetPostgresDBConn(); postgresDBConn != nil {
		logger.Info("shutdown: closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("shutdown: failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	os.Exit(0)
}
