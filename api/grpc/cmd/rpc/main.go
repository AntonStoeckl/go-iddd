package main

import (
	"database/sql"
	"fmt"
	"go-iddd/api/grpc/customer"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/customer/ports/secondary/customers"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"
	"net"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	rpcHostname = "localhost"
	rpcPort     = "5566"
)

var (
	stopSignalChannel chan os.Signal
	logger            *logrus.Logger
	grpcServer        *grpc.Server
	postgresDBConn    *sql.DB
)

func main() {
	buildLogger()
	buildStopSignalChan()
	mustOpenPostgresDBConnection()

	go mustStartGRPC()

	waitForStopSignal()
}

func buildLogger() {
	if logger == nil {
		logger = logrus.New()
		formatter := &logrus.TextFormatter{
			FullTimestamp: true,
		}
		logger.SetFormatter(formatter)
	}
}

func buildStopSignalChan() {
	if stopSignalChannel == nil {
		stopSignalChannel = make(chan os.Signal, 1)
		signal.Notify(stopSignalChannel, os.Interrupt)
	}
}

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		logger.Info("opening Postgres DB handle ...")

		dsn := "postgresql://goiddd:password123@localhost:5432/goiddd_local?sslmode=disable"

		if postgresDBConn, err = sql.Open("postgres", dsn); err != nil {
			logger.Errorf("failed to open Postgres DB handle: %s", err)
			shutdown()
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Errorf("failed to connect to Postgres DB: %s", err)
			shutdown()
		}
	}
}

func mustStartGRPC() {
	logger.Info("starting gRPC server ...")

	rpcHostAndPort := fmt.Sprintf("%s:%s", rpcHostname, rpcPort)

	listener, err := net.Listen("tcp", rpcHostAndPort)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		shutdown()
	}

	grpcServer = grpc.NewServer()
	customerServer := customer.NewCustomerServer(buildCommandHandler())

	customer.RegisterCustomerServer(grpcServer, customerServer)
	reflection.Register(grpcServer)

	logger.Infof("gRPC server ready at %s ...", rpcHostAndPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		shutdown()
	}
}

func buildCommandHandler() shared.CommandHandler {
	es := eventstore.NewPostgresEventStore(postgresDBConn, "eventstore", domain.UnmarshalDomainEvent)
	identityMap := customers.NewIdentityMap()
	repo := customers.NewEventSourcedRepository(es, domain.ReconstituteCustomerFrom, identityMap)
	commandHandler := application.NewCommandHandler(repo, postgresDBConn)

	return commandHandler
}

func waitForStopSignal() {
	s, _ := <-stopSignalChannel

	switch s.(type) {
	case os.Signal:
		logger.Infof("received '%s'", s)
		shutdown()
	}
}

func shutdown() {
	logger.Info("stopping services ...")

	if grpcServer != nil {
		logger.Info("stopping grpc server gracefully ...")
		grpcServer.GracefulStop()
	}

	if postgresDBConn != nil {
		logger.Info("closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("failed to close the Postgres DB connection: %s", err)
		}
	}

	close(stopSignalChannel)

	logger.Info("all services stopped - exiting")

	os.Exit(0)
}
