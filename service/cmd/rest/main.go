package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	customerrest "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/rest"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	var restServer *http.Server
	var grpcClientConn *grpc.ClientConn

	logger := shared.NewStandardLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)

	shutdown := func() {
		shutdown(logger, cancelCtx, grpcClientConn, restServer, func() { os.Exit(1) })
	}

	restServer, grpcClientConn = buildRestServer(config, logger, ctx, shutdown)

	go startRestServer(config, logger, restServer, shutdown)

	waitForStopSignal(logger, shutdown)
}

func buildRestServer(
	config *cmd.Config,
	logger *shared.Logger,
	ctx context.Context,
	shutdown func(),
) (*http.Server, *grpc.ClientConn) {

	logger.Info("configuring REST server ...")

	grpcClientConn, err := grpc.DialContext(ctx, config.GRPC.HostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		shutdown()
	}

	rmux := runtime.NewServeMux(
		runtime.WithProtoErrorHandler(customerrest.CustomHTTPError),
	)

	client := customergrpc.NewCustomerClient(grpcClientConn)

	if err := customerrest.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		shutdown()
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc(
		"/v1/customer/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "service/customeraccounts/infrastructure/adapter/rest/customer.swagger.json")
		},
	)

	restServer := &http.Server{
		Addr:    config.REST.HostAndPort,
		Handler: mux,
	}

	return restServer, grpcClientConn
}

func startRestServer(
	config *cmd.Config,
	logger *shared.Logger,
	restServer *http.Server,
	shutdown func(),
) {

	logger.Infof("starting REST server listening at %s ...", config.REST.HostAndPort)
	logger.Infof("will serve Swagger file at: http://%s/v1/customer/swagger.json", config.REST.HostAndPort)

	if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
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
	cancelCtx context.CancelFunc,
	grpcClientConn *grpc.ClientConn,
	restServer *http.Server,
	exit func(),
) {

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

	exit()
}
