package main

import (
	"context"
	"os"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/service/rest"
	"github.com/AntonStoeckl/go-iddd/src/shared"
)

func main() {
	logger := shared.NewStandardLogger()
	config := rest.MustBuildConfigFromEnv(logger)
	exitFn := func() { os.Exit(1) }
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(config.REST.GRPCDialTimeout)*time.Second)
	grpcClientConn := rest.MustDialGRPCContext(ctx, config, logger, cancelFn)

	s := rest.InitService(ctx, cancelFn, config, logger, exitFn, grpcClientConn)
	go s.StartRestServer()
	s.WaitForStopSignal()
}
