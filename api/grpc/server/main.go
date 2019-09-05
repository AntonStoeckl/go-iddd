package main

import (
	"context"
	"database/sql"
	"go-iddd/api/grpc/customer"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	stopSignalChannel chan os.Signal
	logger            *logrus.Logger
	cancelCtx         context.CancelFunc
	grpcClientConn    *grpc.ClientConn
	grpcServer        *grpc.Server
	postgresDBConn    *sql.DB
)

func main() {
	buildLogger()
	buildStopSignalChan()
	mustOpenPostgresDBConnection()

	go startGRPC()
	go startHTTP()

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
		logger.Info("opening Postgres DB connection ...")

		dsn := "postgresql://goiddd:password123@localhost:5432/goiddd_local?sslmode=disable"

		if postgresDBConn, err = sql.Open("postgres", dsn); err != nil {
			logger.Fatalf("failed to create Postgres DB connection: %s", err)
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Fatalf("failed to connect to Postgres DB: %s", err)
		}
	}
}

//func foobar() {
//	es := eventstore.NewPostgresEventStore(postgresDBConn, "eventstore", domain.UnmarshalDomainEvent)
//	identityMap := customers.NewIdentityMap()
//	repo := customers.NewEventSourcedRepository(es, domain.ReconstituteCustomerFrom, identityMap)
//	application.NewCommandHandler(repo, postgresDBConn)
//}

func startGRPC() {
	logger.Info("starting gRPC server ...")

	listener, err := net.Listen("tcp", "localhost:5566")
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		stopSignalChannel <- os.Interrupt
	}

	grpcServer = grpc.NewServer()
	customerServer := customer.NewCustomerServer()

	customer.RegisterCustomerServer(grpcServer, customerServer)
	reflection.Register(grpcServer)

	logger.Info("gRPC server ready ...")

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		stopSignalChannel <- os.Interrupt
	}
}

func startHTTP() {
	var err error
	var ctx context.Context

	logger.Info("starting REST server ...")

	ctx, cancelCtx = context.WithCancel(context.Background())

	grpcClientConn, err = grpc.Dial("localhost:5566", grpc.WithInsecure())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		stopSignalChannel <- os.Interrupt
	}

	rmux := runtime.NewServeMux()
	client := customer.NewCustomerClient(grpcClientConn)

	if err = customer.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		stopSignalChannel <- os.Interrupt
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	serveSwagger := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/swagger.json")
	}

	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("www/swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	logger.Info("REST server ready - serving Swagger at: http://localhost:8080/swagger-ui/")

	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
		stopSignalChannel <- os.Interrupt
	}
}

func waitForStopSignal() {
	s := <-stopSignalChannel

	logger.Infof("received '%s' - stopping services ...", s)

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
}
