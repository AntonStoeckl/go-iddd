package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/primary/grpc"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ctxTimeout = 10 * time.Second
)

var (
	diContainer    *cmd.DIContainer
	grpcServer     *grpc.Server
	grpcClientConn *grpc.ClientConn
	cancelCtx      context.CancelFunc
	restServer     *http.Server
)

func main() {
	var err error

	logger := cmd.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)

	diContainer, err = cmd.Bootstrap(config, logger)
	if err != nil {
		shutdown(logger)
	}

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	go mustStartGRPC(config, logger)
	go mustStartREST(config, logger)

	waitForStopSignal(stopSignalChannel, logger)
}

func mustStartGRPC(config *cmd.Config, logger *cmd.Logger) {
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

func mustStartREST(config *cmd.Config, logger *cmd.Logger) {
	var err error
	var ctx context.Context

	logger.Info("configuring REST server ...")

	ctx, cancelCtx = context.WithTimeout(context.Background(), ctxTimeout)

	grpcClientConn, err = grpc.DialContext(ctx, config.GRPC.HostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		shutdown(logger)
	}

	rmux := runtime.NewServeMux(
		runtime.WithProtoErrorHandler(customergrpc.CustomHTTPError),
	)

	client := customergrpc.NewCustomerClient(grpcClientConn)

	if err = customergrpc.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		shutdown(logger)
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

	restServer = &http.Server{
		Addr:    config.REST.HostAndPort,
		Handler: mux,
	}

	logger.Infof("starting REST server listening at %s ...", config.REST.HostAndPort)
	logger.Infof("will serve Swagger file at: http://%s/v1/customer/swagger.json", config.REST.HostAndPort)

	if err = restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
		shutdown(logger)
	}
}

func waitForStopSignal(stopSignalChannel chan os.Signal, logger *cmd.Logger) {
	logger.Info("start waiting for stop signal ...")

	s := <-stopSignalChannel

	switch s.(type) {
	case os.Signal:
		logger.Infof("received '%s'", s)
		close(stopSignalChannel)
		shutdown(logger)
	}
}

func shutdown(logger *cmd.Logger) {
	logger.Info("shutdown: stopping services ...")

	if cancelCtx != nil {
		logger.Info("shutdown: canceling context ...")
		cancelCtx()
	}

	if restServer != nil {
		logger.Info("shutdown: stopping REST server gracefully ...")
		if err := restServer.Shutdown(context.Background()); err != nil {
			logger.Warnf("shutdown: failed to stop the REST server: %s", err)
		}
	}

	if grpcClientConn != nil {
		logger.Info("shutdown: closing gRPC client connection ...")

		if err := grpcClientConn.Close(); err != nil {
			logger.Warnf("shutdown: failed to close the gRPC client connection: %s", err)
		}
	}

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
