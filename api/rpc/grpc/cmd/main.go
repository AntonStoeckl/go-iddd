package main

import (
	"go-iddd/api/rpc/grpc/customer"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go startGRPC()

	// Block forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
