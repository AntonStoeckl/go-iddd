package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	customergrpcproto "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc/proto"
	customerrest "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/rest"
	customerrestproto "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/rest/proto"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	config         *Config
	logger         *shared.Logger
	exitFn         func()
	ctx            context.Context
	cancelFn       context.CancelFunc
	restServer     *http.Server
	grpcClientConn *grpc.ClientConn
}

func MustDialGRPCContext(
	ctx context.Context,
	config *Config,
	logger *shared.Logger,
	cancelFn context.CancelFunc,
) *grpc.ClientConn {

	grpcClientConn, err := grpc.DialContext(ctx, config.REST.GRPCDialHostAndPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancelFn()
		logger.Panic().Msgf("fail to dial gRPC service: %s", err)
	}

	return grpcClientConn
}

func InitService(
	ctx context.Context,
	cancelFn context.CancelFunc,
	config *Config,
	logger *shared.Logger,
	exitFn func(),
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
	s.logger.Info().Msg("configuring REST server ...")

	client := customergrpcproto.NewCustomerClient(s.grpcClientConn)

	rmux := runtime.NewServeMux(
		runtime.WithProtoErrorHandler(customerrest.CustomHTTPError),
	)

	if err := customerrestproto.RegisterCustomerHandlerClient(s.ctx, rmux, client); err != nil {
		s.logger.Error().Msgf("failed to register customerHandlerClient: %s", err)
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
	s.logger.Info().Msgf("starting REST server listening at %s ...", hostAndPort)
	s.logger.Info().Msgf("will serve Swagger file at: http://%s/v1/customer/swagger.json", hostAndPort)

	if err := s.restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error().Msgf("REST server failed to listenAndServe: %s", err)
		s.shutdown()
	}
}

func (s *Service) WaitForStopSignal() {
	s.logger.Info().Msg("start waiting for stop signal ...")

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-stopSignalChannel

	if _, ok := sig.(os.Signal); ok {
		s.logger.Info().Msgf("received '%s'", sig)
		close(stopSignalChannel)
		s.shutdown()
	}
}

func (s *Service) shutdown() {
	s.logger.Info().Msg("shutdown: stopping services ...")

	if s.cancelFn != nil {
		s.logger.Info().Msg("shutdown: canceling context ...")
		s.cancelFn()
	}

	if s.restServer != nil {
		s.logger.Info().Msg("shutdown: stopping REST server gracefully ...")
		if err := s.restServer.Shutdown(context.Background()); err != nil {
			s.logger.Warn().Msgf("shutdown: failed to stop the REST server: %s", err)
		}
	}

	if s.grpcClientConn != nil {
		s.logger.Info().Msg("shutdown: closing gRPC client connection ...")
		if err := s.grpcClientConn.Close(); err != nil {
			s.logger.Warn().Msgf("shutdown: failed to close the gRPC client connection: %s", err)
		}
	}

	s.logger.Info().Msg("shutdown: all services stopped - Hasta la vista, baby!")

	s.exitFn()
}
