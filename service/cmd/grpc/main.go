package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/primary/grpc"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	diContainer *cmd.DIContainer
	grpcServer  *grpc.Server
)

func main() {
	var err error

	logger := shared.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)

	diContainer, err = cmd.Bootstrap(config, logger)
	if err != nil {
		shutdown(logger)
	}

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	go mustStartGRPC(config, logger)

	waitForStopSignal(stopSignalChannel, logger)
}

func mustStartGRPC(config *cmd.Config, logger *shared.Logger) {
	logger.Info("configuring gRPC server ...")

	listener, err := net.Listen("tcp", config.GRPC.HostAndPort)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		shutdown(logger)
	}

	grpcServer = grpc.NewServer()
	customergrpc.RegisterCustomerServer(grpcServer, diContainer.GetCustomerGRPCServer())
	reflection.Register(grpcServer)

	logger.Infof("starting gRPC server listening at %s ...", config.GRPC.HostAndPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		shutdown(logger)
	}
}

func waitForStopSignal(stopSignalChannel chan os.Signal, logger *shared.Logger) {
	logger.Info("start waiting for stop signal ...")

	s := <-stopSignalChannel

	switch s.(type) {
	case os.Signal:
		logger.Infof("received '%s'", s)
		close(stopSignalChannel)
		shutdown(logger)
	}
}

func shutdown(logger *shared.Logger) {
	logger.Info("shutdown: stopping services ...")

	if grpcServer != nil {
		logger.Info("shutdown: stopping gRPC server gracefully ...")
		grpcServer.GracefulStop()
	}

	postgresDBConn := diContainer.GetPostgresDBConn()

	if postgresDBConn != nil {
		logger.Info("shutdown: closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("shutdown: failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	os.Exit(0)
}
