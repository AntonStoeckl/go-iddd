package grpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntonStoeckl/go-iddd/src/service"
	"github.com/AntonStoeckl/go-iddd/src/shared"
)

type Service struct {
	config       *service.Config
	logger       *shared.Logger
	diContainter *service.DIContainer
	exitFn       func()
}

func InitService(
	config *service.Config,
	logger *shared.Logger,
	exitFn func(),
	diContainter *service.DIContainer,
) *Service {

	return &Service{
		config:       config,
		logger:       logger,
		exitFn:       exitFn,
		diContainter: diContainter,
	}
}

func (s Service) StartGRPCServer() {
	s.logger.Info("configuring gRPC server ...")

	listener, err := net.Listen("tcp", s.config.GRPC.HostAndPort)
	if err != nil {
		s.logger.Errorf("failed to listen: %v", err)
		s.shutdown()
	}

	s.logger.Infof("starting gRPC server listening at %s ...", s.config.GRPC.HostAndPort)

	grpcServer := s.diContainter.GetGRPCServer()
	if err := grpcServer.Serve(listener); err != nil {
		s.logger.Errorf("gRPC server failed to serve: %s", err)
		s.shutdown()
	}
}

func (s Service) WaitForStopSignal() {
	s.logger.Info("start waiting for stop signal ...")

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-stopSignalChannel

	switch sig.(type) {
	case os.Signal:
		s.logger.Infof("received '%s'", sig)
		close(stopSignalChannel)
		s.shutdown()
	}
}

func (s Service) shutdown() {
	s.logger.Info("shutdown: stopping services ...")

	grpcServer := s.diContainter.GetGRPCServer()
	if grpcServer != nil {
		s.logger.Info("shutdown: stopping gRPC server gracefully ...")
		grpcServer.GracefulStop()
	}

	postgresDBConn := s.diContainter.GetPostgresDBConn()
	if postgresDBConn != nil {
		s.logger.Info("shutdown: closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			s.logger.Warnf("shutdown: failed to close the Postgres DB connection: %s", err)
		}
	}

	s.logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	s.exitFn()
}
