package main

import (
	"os"

	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	stdLogger := shared.NewStandardLogger()
	config := grpc.MustBuildConfigFromEnv(stdLogger)
	exitFn := func() { os.Exit(1) }
	// @TODO check nil reference error if db connection be nil
	postgresDBConn := grpc.MustInitPostgresDB(config, stdLogger)
	mongoDBConn := grpc.MustInitMongoDB(config, stdLogger)
	diContainer := grpc.MustBuildDIContainer(
		config,
		stdLogger,
		grpc.UsePostgresDBConn(postgresDBConn),
		grpc.UseMongoDBConn(mongoDBConn),
	)

	s := grpc.InitService(config, stdLogger, exitFn, diContainer)
	go s.StartGRPCServer()
	s.WaitForStopSignal()
}
