package main

import (
	"context"
	"go-iddd/api/grpc/customer"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	signalChan     chan os.Signal
	logger         *logrus.Logger
	cancelCtx      context.CancelFunc
	grpcClientConn *grpc.ClientConn
)

func main() {
	createSignalChan()
	buildLogger()

	go startGRPC()
	go startHTTP()

	waitUntilStopped()
}

func buildLogger() {
	logger = logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.SetFormatter(formatter)
}

func startGRPC() {
	listener, err := net.Listen("tcp", "localhost:5566")
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		signalChan <- os.Interrupt
	}

	grpcServer := grpc.NewServer()
	customerServer := customer.NewCustomerServer()

	customer.RegisterCustomerServer(grpcServer, customerServer)
	reflection.Register(grpcServer)

	logger.Info("gRPC server ready ...")

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		signalChan <- os.Interrupt
	}
}

func startHTTP() {
	var err error
	var ctx context.Context

	ctx, cancelCtx = context.WithCancel(context.Background())

	grpcClientConn, err = grpc.Dial("localhost:5566", grpc.WithInsecure())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		signalChan <- os.Interrupt
	}

	rmux := runtime.NewServeMux()
	client := customer.NewCustomerClient(grpcClientConn)

	if err = customer.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		signalChan <- os.Interrupt
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("www/swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	logger.Info("REST server ready ...")
	logger.Info("Serving Swagger at: http://localhost:8080/swagger-ui/")

	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
		signalChan <- os.Interrupt
	}
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/swagger.json")
}

func createSignalChan() {
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
}

func waitUntilStopped() {
	s := <-signalChan

	logger.Infof("received '%s' - stopping services ...\n", s)

	if cancelCtx != nil {
		logger.Info("canceling context ...")
		cancelCtx()
	}

	if grpcClientConn != nil {
		logger.Info("closing grpc client connection ...")

		if err := grpcClientConn.Close(); err != nil {
			logger.Warnf("failed to close the grpc client connection: %s", err)
		}
	}

	close(signalChan)

	logger.Info("all services stopped - exiting")
}
