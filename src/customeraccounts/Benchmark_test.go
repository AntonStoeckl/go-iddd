package customeraccounts_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
)

type benchmarkTestValues struct {
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
	config := grpc.MustBuildConfigFromEnv(logger)
	postgresDBConn := grpc.MustInitPostgresDB(config, logger)
	diContainer := grpc.MustBuildDIContainer(config, logger, grpc.UsePostgresDBConn(postgresDBConn))
	commandHandler := diContainer.GetCustomerCommandHandler()
	v := initBenchmarkTestValues()
	prepareForBenchmark(b, commandHandler, &v)

	b.Run("ChangeName", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if n%2 == 0 {
				if err = commandHandler.ChangeCustomerName(v.customerID.String(), v.newGivenName, v.newFamilyName); err != nil {
					b.FailNow()
				}
			} else {
				if err = commandHandler.ChangeCustomerName(v.customerID.String(), v.givenName, v.familyName); err != nil {
					b.FailNow()
				}
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		commandHandler,
		v.customerID,
	)
}

func BenchmarkCustomerQuery(b *testing.B) {
	logger := shared.NewNilLogger()
	config := grpc.MustBuildConfigFromEnv(logger)
	postgresDBConn := grpc.MustInitPostgresDB(config, logger)
	diContainer := grpc.MustBuildDIContainer(config, logger, grpc.UsePostgresDBConn(postgresDBConn))
	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	v := initBenchmarkTestValues()
	prepareForBenchmark(b, commandHandler, &v)

	b.Run("CustomerViewByID", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err := queryHandler.CustomerViewByID(v.customerID.String()); err != nil {
				b.FailNow()
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		commandHandler,
		v.customerID,
	)
}

func initBenchmarkTestValues() benchmarkTestValues {
	var v benchmarkTestValues

	v.emailAddress = "fiona@gallagher.net"
	v.givenName = "Fiona"
	v.familyName = "Galagher"
	v.newEmailAddress = "fiona@pratt.net"
	v.newGivenName = "Fiona"
	v.newFamilyName = "Pratt"

	return v
}

func prepareForBenchmark(
	b *testing.B,
	commandHandler *application.CustomerCommandHandler,
	v *benchmarkTestValues,
) {

	var err error

	v.customerID = value.GenerateCustomerID()

	if err = commandHandler.RegisterCustomer(v.customerID, v.emailAddress, v.givenName, v.familyName); err != nil {
		b.FailNow()
	}

	for n := 0; n < 100; n++ {
		if n%2 == 0 {
			if err = commandHandler.ChangeCustomerEmailAddress(v.customerID.String(), v.newEmailAddress); err != nil {
				b.FailNow()
			}
		} else {
			if err = commandHandler.ChangeCustomerEmailAddress(v.customerID.String(), v.emailAddress); err != nil {
				b.FailNow()
			}
		}
	}
}

func cleanUpAfterBenchmark(
	b *testing.B,
	eventstore application.EventStoreInterface,
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
