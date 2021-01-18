package main

import (
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
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
		cmd.UsePostgresDBConn(postgresDBConn),
	)
	grpcServer := diContainer.GetGRPCServer()

	shutdown := func() {
		shutdown(logger, grpcServer, postgresDBConn, func() { os.Exit(1) })
	}

	go startGRPCServer(config, logger, grpcServer, shutdown)

	waitForStopSignal(logger, shutdown)
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

func waitForStopSignal(logger *shared.Logger, shutdown func()) {
	logger.Info("start waiting for stop signal ...")

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-stopSignalChannel

	switch sig.(type) {
	case os.Signal:
		logger.Infof("received '%s'", sig)
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
