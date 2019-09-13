package main

import (
	"database/sql"
	"fmt"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/di"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	logger         *logrus.Logger
	postgresDBConn *sql.DB
	diContainer    *di.Container
)

func main() {
	bootstrap()
	mustRunCLIApp()
}

func mustRunCLIApp() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:      "RegisterCustomer",
			Aliases:   []string{"rc"},
			Usage:     "Register a Customer",
			Action:    registerCustomer,
			ArgsUsage: "emailAddress givenName familyName",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func registerCustomer(ctx *cli.Context) error {
	emailAddress := ctx.Args().Get(0)
	givenName := ctx.Args().Get(1)
	familyName := ctx.Args().Get(2)
	id := values.GenerateID()

	commandHandler := diContainer.GetCustomerCommandHandler()

	register, err := commands.NewRegister(id.String(), emailAddress, givenName, familyName)
	if err != nil {
		return err
	}

	if err := commandHandler.Handle(register); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(ctx.App.Writer, "customer registered with id: %s\n", id.String())

	return nil
}

func bootstrap() {
	buildLogger()
	mustOpenPostgresDBConnection()
	mustBuildDIContainer()
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

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		//logger.Info("opening Postgres DB handle ...")

		dsn := "postgresql://goiddd:password123@localhost:5432/goiddd_local?sslmode=disable"

		if postgresDBConn, err = sql.Open("postgres", dsn); err != nil {
			logger.Errorf("failed to open Postgres DB handle: %s", err)
			shutdown()
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Errorf("failed to connect to Postgres DB: %s", err)
			shutdown()
		}
	}
}

func mustBuildDIContainer() {
	var err error

	if diContainer == nil {
		if diContainer, err = di.NewContainer(postgresDBConn); err != nil {
			logger.Errorf("failed to build the DI container: %s", err)
			shutdown()
		}
	}
}

func shutdown() {
	logger.Info("stopping services ...")

	if postgresDBConn != nil {
		//logger.Info("closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("all services stopped - exiting")

	os.Exit(0)
}
