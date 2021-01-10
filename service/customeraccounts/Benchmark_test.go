package customeraccounts_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/serialization"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/postgres"
	"github.com/AntonStoeckl/go-iddd/service/shared"
)

type benchmarkTestArtifacts struct {
	customerID      value.CustomerID
	emailAddress    string
	givenName       string
	familyName      string
	newEmailAddress string
	newGivenName    string
	newFamilyName   string
}

func BenchmarkCustomerCommand(b *testing.B) {
	var err error

	logger := shared.NewNilLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	postgresDBConn := cmd.MustInitPostgresDB(config, logger)
	diContainer := cmd.MustBuildDIContainer(config, logger,
		serialization.MarshalCustomerEvent,
		serialization.UnmarshalCustomerEvent,
		customer.BuildUniqueEmailAddressAssertions,
		cmd.WithPostgresDBConn(postgresDBConn),
	)

	commandHandler := diContainer.GetCustomerCommandHandler()
	ba := buildArtifactsForBenchmarkTest()
	prepareForBenchmark(b, commandHandler, &ba)

	b.Run("ChangeName", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if n%2 == 0 {
				if err = commandHandler.ChangeCustomerName(ba.customerID.String(), ba.newGivenName, ba.newFamilyName); err != nil {
					b.FailNow()
				}
			} else {
				if err = commandHandler.ChangeCustomerName(ba.customerID.String(), ba.givenName, ba.familyName); err != nil {
					b.FailNow()
				}
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		commandHandler,
		ba.customerID,
	)
}

func BenchmarkCustomerQuery(b *testing.B) {
	logger := shared.NewNilLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	postgresDBConn := cmd.MustInitPostgresDB(config, logger)
	diContainer := cmd.MustBuildDIContainer(config, logger,
		serialization.MarshalCustomerEvent,
		serialization.UnmarshalCustomerEvent,
		customer.BuildUniqueEmailAddressAssertions,
		cmd.WithPostgresDBConn(postgresDBConn),
	)

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	ba := buildArtifactsForBenchmarkTest()
	prepareForBenchmark(b, commandHandler, &ba)

	b.Run("CustomerViewByID", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err := queryHandler.CustomerViewByID(ba.customerID.String()); err != nil {
				b.FailNow()
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		commandHandler,
		ba.customerID,
	)
}

func buildArtifactsForBenchmarkTest() benchmarkTestArtifacts {
	var ba benchmarkTestArtifacts

	ba.emailAddress = "fiona@gallagher.net"
	ba.givenName = "Fiona"
	ba.familyName = "Galagher"
	ba.newEmailAddress = "fiona@pratt.net"
	ba.newGivenName = "Fiona"
	ba.newFamilyName = "Pratt"

	return ba
}

func prepareForBenchmark(
	b *testing.B,
	commandHandler *application.CustomerCommandHandler,
	ba *benchmarkTestArtifacts,
) {

	var err error

	if ba.customerID, err = commandHandler.RegisterCustomer(ba.emailAddress, ba.givenName, ba.familyName); err != nil {
		b.FailNow()
	}

	for n := 0; n < 100; n++ {
		if n%2 == 0 {
			if err = commandHandler.ChangeCustomerEmailAddress(ba.customerID.String(), ba.newEmailAddress); err != nil {
				b.FailNow()
			}
		} else {
			if err = commandHandler.ChangeCustomerEmailAddress(ba.customerID.String(), ba.emailAddress); err != nil {
				b.FailNow()
			}
		}
	}
}

func cleanUpAfterBenchmark(
	b *testing.B,
	eventstore *postgres.CustomerEventStore,
	commandHandler *application.CustomerCommandHandler,
	id value.CustomerID,
) {

	if err := commandHandler.DeleteCustomer(id.String()); err != nil {
		b.FailNow()
	}

	if err := eventstore.PurgeEventStream(id); err != nil {
		b.FailNow()
	}
}
