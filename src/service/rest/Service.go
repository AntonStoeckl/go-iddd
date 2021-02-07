package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	customergrpc "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc"
	customerrest "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/rest"
	"github.com/AntonStoeckl/go-iddd/src/service"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	config         *service.Config
	logger         *shared.Logger
	exitFn         func()
	ctx            context.Context
	cancelFn       context.CancelFunc
	restServer     *http.Server
	grpcClientConn *grpc.ClientConn
}

func MustDialGRPCContext(
	config *service.Config,
	logger *shared.Logger,
	ctx context.Context,
	cancelFn context.CancelFunc,
) *grpc.ClientConn {

	grpcClientConn, err := grpc.DialContext(ctx, config.GRPC.HostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancelFn()
		logger.Panicf("fail to dial gRPC service: %s", err)
	}

	return grpcClientConn
}

func InitService(
	config *service.Config,
	logger *shared.Logger,
	exitFn func(),
	ctx context.Context,
	cancelFn context.CancelFunc,
	grpcClientConn *grpc.ClientConn,

) *Service {

	s := &Service{
		config:         config,
		logger:         logger,
		exitFn:         exitFn,
		ctx:            ctx,
		cancelFn:       cancelFn,
		grpcClientConn: grpcClientConn,
	}

	s.buildRestServer()

	return s
}

func (s *Service) buildRestServer() {
	s.logger.Info("configuring REST server ...")

	client := customergrpc.NewCustomerClient(s.grpcClientConn)

	rmux := runtime.NewServeMux(
		runtime.WithProtoErrorHandler(customerrest.CustomHTTPError),
	)

	if err := customerrest.RegisterCustomerHandlerClient(s.ctx, rmux, client); err != nil {
		s.logger.Errorf("failed to register customerHandlerClient: %s", err)
		s.shutdown()
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	// Serve the swagger file and swagger-ui (really?)
	mux.HandleFunc(
		"/v1/customer/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, fmt.Sprintf("%s/customer.swagger.json", s.config.REST.SwaggerFilePathCustomer))
		},
	)

	s.restServer = &http.Server{
		Addr:    s.config.REST.HostAndPort,
		Handler: mux,
	}
}

func (s *Service) StartRestServer() {
	hostAndPort := s.config.REST.HostAndPort
	s.logger.Infof("starting REST server listening at %s ...", hostAndPort)
	s.logger.Infof("will serve Swagger file at: http://%s/v1/customer/swagger.json", hostAndPort)

	if err := s.restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Errorf("REST server failed to listenAndServe: %s", err)
		s.shutdown()
	}
}

func (s *Service) WaitForStopSignal() {
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

func (s *Service) shutdown() {
	s.logger.Info("shutdown: stopping services ...")

	if s.cancelFn != nil {
		s.logger.Info("shutdown: canceling context ...")
		s.cancelFn()
	}

	if s.restServer != nil {
		s.logger.Info("shutdown: stopping REST server gracefully ...")
		if err := s.restServer.Shutdown(context.Background()); err != nil {
			s.logger.Warnf("shutdown: failed to stop the REST server: %s", err)
		}
	}

	if s.grpcClientConn != nil {
		s.logger.Info("shutdown: closing gRPC client connection ...")
		if err := s.grpcClientConn.Close(); err != nil {
			s.logger.Warnf("shutdown: failed to close the gRPC client connection: %s", err)
		}
	}

	s.logger.Info("shutdown: all services stopped - Hasta la vista, baby!")

	s.exitFn()
}
