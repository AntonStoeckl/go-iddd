package main

import (
	"go-iddd/api/grpc/customer"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	signalChan chan os.Signal
)

func main() {
	createSignalChan()

	go startGRPC()

	waitUntilStopped()
}

func startGRPC() {
	lis, err := net.Listen("tcp", "localhost:5566")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	customerServer := customer.NewCustomerServer()

	customer.RegisterCustomerServer(grpcServer, customerServer)

	reflection.Register(grpcServer)

	log.Println("gRPC server ready...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func createSignalChan() {
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
}

func waitUntilStopped() {
	<-signalChan
	log.Println("service stopped")
}
