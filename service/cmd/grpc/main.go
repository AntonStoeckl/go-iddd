package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-iddd/customer/domain/events"
	customergrpc "go-iddd/customer/infrastructure/grpc"
	"go-iddd/service"
	"go-iddd/shared/infrastructure/eventstore"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	rpcHostname  = "localhost"
	rpcPort      = "5566"
	ctxTimeout   = 10 * time.Second
	restHostname = "localhost"
	restPort     = "8080"
)

var (
	config            *service.Config
	logger            *service.Logger
	postgresDBConn    *sql.DB
	diContainer       *service.DIContainer
	grpcServer        *grpc.Server
	cancelCtx         context.CancelFunc
	grpcClientConn    *grpc.ClientConn
	restServer        *http.Server
	stopSignalChannel chan os.Signal
)

func main() {
	bootstrap()
	go mustStartGRPC()
	go mustStartREST()
	waitForStopSignal()
}

func bootstrap() {
	mustBuildConfig()
	buildLogger()
	mustOpenPostgresDBConnection()
	mustRunDBMigrations()
	mustBuildDIContainer()
	buildStopSignalChan()
}

func mustBuildConfig() {
	if config == nil {
		var err error

		config, err = service.NewConfigFromEnv()
		if err != nil {
			logrus.Fatalf("failed to get config from env - exiting: %s", err)
		}
	}
}

func buildLogger() {
	if logger == nil {
		logger = service.NewStandardLogger()
	}
}

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		logger.Info("opening Postgres DB connection ...")

		if postgresDBConn, err = sql.Open("postgres", config.Postgres.DSN); err != nil {
			logger.Errorf("failed to open Postgres DB connection: %s", err)
			shutdown()
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Errorf("failed to connect to Postgres DB: %s", err)
			shutdown()
		}
	}
}

func mustRunDBMigrations() {
	migrator, err := eventstore.NewMigrator(postgresDBConn, config.Postgres.MigrationsPath)
	if err != nil {
		logger.Errorf("failed to create DB migrator: %s", err)
		shutdown()
		return
	}

	err = migrator.WithLogger(logger).Up()
	if err != nil {
		logger.Errorf("failed to run DB migrator: %s", err)
		shutdown()
	}
}

func mustBuildDIContainer() {
	var err error

	if diContainer == nil {
		diContainer, err = service.NewDIContainer(
			postgresDBConn,
			events.UnmarshalDomainEvent,
		)

		if err != nil {
			logger.Errorf("failed to build the DI container: %s", err)
			shutdown()
		}
	}
}

func buildStopSignalChan() {
	if stopSignalChannel == nil {
		stopSignalChannel = make(chan os.Signal, 1)
		signal.Notify(stopSignalChannel, os.Interrupt)
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
	customerServer := customergrpc.NewCustomerServer(diContainer.GetCustomerCommandHandler())

	customergrpc.RegisterCustomerServer(grpcServer, customerServer)
	reflection.Register(grpcServer)

	logger.Infof("gRPC server ready at %s ...", rpcHostAndPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		shutdown()
	}
}

func mustStartREST() {
	var err error
	var ctx context.Context

	logger.Info("starting REST server ...")

	ctx, cancelCtx = context.WithTimeout(context.Background(), ctxTimeout)

	rpcHostAndPort := fmt.Sprintf("%s:%s", rpcHostname, rpcPort)

	grpcClientConn, err = grpc.DialContext(ctx, rpcHostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		shutdown()
	}

	rmux := runtime.NewServeMux()
	client := customergrpc.NewCustomerClient(grpcClientConn)

	if err = customergrpc.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		shutdown()
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc(
		"/v1/customer/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "customer/infrastructure/grpc/customer.swagger.json")
		},
	)

	restHostAndPort := fmt.Sprintf("%s:%s", restHostname, restPort)

	restServer = &http.Server{
		Addr:    restHostAndPort,
		Handler: mux,
	}

	logger.Info("REST server ready")
	logger.Infof("Serving Swagger file at: http://%s/v1/customer/swagger.json", restHostAndPort)

	if err = restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
		shutdown()
	}
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

	if cancelCtx != nil {
		logger.Info("canceling context ...")
		cancelCtx()
	}

	if restServer != nil {
		logger.Info("stopping rest server gracefully ...")
		if err := restServer.Shutdown(context.Background()); err != nil {
			logger.Warnf("failed to stop the rest server: %s", err)
		}
	}

	if grpcClientConn != nil {
		logger.Info("closing grpc client connection ...")

		if err := grpcClientConn.Close(); err != nil {
			logger.Warnf("failed to close the grpc client connection: %s", err)
		}
	}

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
