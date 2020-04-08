package application_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
)

type benchmarkTestArtifacts struct {
	customerID      values.CustomerID
	emailAddress    string
	givenName       string
	familyName      string
	newEmailAddress string
	newGivenName    string
	newFamilyName   string
}

func BenchmarkCustomerCommand(b *testing.B) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

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
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

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
	commandHandler *command.CustomerCommandHandler,
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
	eventstore *eventstore.CustomerEventStore,
	commandHandler *command.CustomerCommandHandler,
	id values.CustomerID,
) {

	if err := commandHandler.DeleteCustomer(id.String()); err != nil {
		b.FailNow()
	}

	if err := eventstore.Purge(id); err != nil {
		b.FailNow()
	}
}
