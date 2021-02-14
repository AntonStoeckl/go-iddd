package main

import (
	"os"

	"github.com/AntonStoeckl/go-iddd/src/service"
	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	stdLogger := shared.NewStandardLogger()
	config := service.MustBuildConfigFromEnv(stdLogger)
	exitFn := func() { os.Exit(1) }
	postgresDBConn := grpc.MustInitPostgresDB(config, stdLogger)
	diContainer := grpc.MustBuildDIContainer(
		config,
		stdLogger,
		grpc.UsePostgresDBConn(postgresDBConn),
	)

	s := grpc.InitService(config, stdLogger, exitFn, diContainer)
	go s.StartGRPCServer()
	s.WaitForStopSignal()
}
