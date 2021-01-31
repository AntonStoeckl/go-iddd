package main

import (
	"context"
	"os"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/service"
	"github.com/AntonStoeckl/go-iddd/src/service/rest"
	"github.com/AntonStoeckl/go-iddd/src/shared"
)

func main() {
	logger := shared.NewStandardLogger()
	config := service.MustBuildConfigFromEnv(logger)
	exitFn := func() { os.Exit(1) }
	ctx, cancelFn := context.WithTimeout(context.Background(), 3*time.Second)
	grpcClientConn := rest.MustDialGRPCContext(config, logger, ctx, cancelFn)

	s := rest.InitService(config, logger, exitFn, ctx, cancelFn, grpcClientConn)
	go s.StartRestServer()
	s.WaitForStopSignal()
}
