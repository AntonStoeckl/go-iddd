package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/primary/grpc"
	customerrest "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/primary/rest"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	ctxTimeout = 3 * time.Second
)

var (
	grpcClientConn *grpc.ClientConn
	cancelCtx      context.CancelFunc
	restServer     *http.Server
)

func main() {
	logger := shared.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt)

	go mustStartREST(config, logger)

	waitForStopSignal(stopSignalChannel, logger)
}

func mustStartREST(config *cmd.Config, logger *shared.Logger) {
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
		runtime.WithProtoErrorHandler(customerrest.CustomHTTPError),
	)

	client := customergrpc.NewCustomerClient(grpcClientConn)

	if err = customerrest.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		shutdown(logger)
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc(
		"/v1/customer/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "service/customeraccounts/infrastructure/adapter/primary/rest/customer.swagger.json")
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

	logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	os.Exit(0)
}
