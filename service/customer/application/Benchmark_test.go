package application_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
)

type benchmarkTestArtifacts struct {
	registerCustomer               commands.RegisterCustomer
	confirmCustomerEmailAddress    commands.ConfirmCustomerEmailAddress
	changeCustomerEmailAddress     commands.ChangeCustomerEmailAddress
	changeCustomerEmailAddressBack commands.ChangeCustomerEmailAddress
	changeCustomerName             commands.ChangeCustomerName
	changeCustomerNameBack         commands.ChangeCustomerName
}

func BenchmarkCustomerCommand(b *testing.B) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	ba := buildArtifactsForBenchmarkTest(b)
	prepareForBenchmark(b, commandHandler, ba)

	b.Run("ChangeName", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if n%2 == 0 {
				if err = commandHandler.ChangeCustomerName(ba.changeCustomerName); err != nil {
					b.FailNow()
				}
			} else {
				if err = commandHandler.ChangeCustomerName(ba.changeCustomerNameBack); err != nil {
					b.FailNow()
				}
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		ba.registerCustomer.CustomerID(),
	)
}

func BenchmarkCustomerQuery(b *testing.B) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	ba := buildArtifactsForBenchmarkTest(b)
	prepareForBenchmark(b, commandHandler, ba)

	b.Run("CustomerViewByID", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err := queryHandler.CustomerViewByID(ba.registerCustomer.CustomerID()); err != nil {
				b.FailNow()
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		ba.registerCustomer.CustomerID(),
	)
}

func buildArtifactsForBenchmarkTest(b *testing.B) benchmarkTestArtifacts {
	var err error
	var ba benchmarkTestArtifacts

	emailAddress := "fiona@gallagher.net"
	givenName := "Fiona"
	familyName := "Galagher"
	newEmailAddress := "fiona@pratt.net"
	newGivenName := "Fiona"
	newFamilyName := "Pratt"

	ba.registerCustomer, err = commands.BuildRegisterCustomer(
		emailAddress,
		givenName,
		familyName,
	)

	if err != nil {
		b.FailNow()
	}

	ba.confirmCustomerEmailAddress, err = commands.BuildConfirmCustomerEmailAddress(
		ba.registerCustomer.CustomerID().ID(),
		ba.registerCustomer.ConfirmationHash().Hash(),
	)

	if err != nil {
		b.FailNow()
	}

	ba.changeCustomerEmailAddress, err = commands.BuildChangeCustomerEmailAddress(
		ba.registerCustomer.CustomerID().ID(),
		newEmailAddress,
	)

	if err != nil {
		b.FailNow()
	}

	ba.changeCustomerEmailAddressBack, err = commands.BuildChangeCustomerEmailAddress(
		ba.registerCustomer.CustomerID().ID(),
		emailAddress,
	)

	if err != nil {
		b.FailNow()
	}

	ba.changeCustomerName, err = commands.BuildChangeCustomerName(
		ba.registerCustomer.CustomerID().ID(),
		newGivenName,
		newFamilyName,
	)

	if err != nil {
		b.FailNow()
	}

	ba.changeCustomerNameBack, err = commands.BuildChangeCustomerName(
		ba.registerCustomer.CustomerID().ID(),
		givenName,
		familyName,
	)

	if err != nil {
		b.FailNow()
	}

	return ba
}

func prepareForBenchmark(
	b *testing.B,
	commandHandler *command.CustomerCommandHandler,
	ba benchmarkTestArtifacts,
) {

	var err error

	if err = commandHandler.RegisterCustomer(ba.registerCustomer); err != nil {
		b.FailNow()
	}

	if err = commandHandler.ConfirmCustomerEmailAddress(ba.confirmCustomerEmailAddress); err != nil {
		b.FailNow()
	}

	for n := 0; n < 100; n++ {
		if n%2 == 0 {
			if err = commandHandler.ChangeCustomerEmailAddress(ba.changeCustomerEmailAddress); err != nil {
				b.FailNow()
			}
		} else {
			if err = commandHandler.ChangeCustomerEmailAddress(ba.changeCustomerEmailAddressBack); err != nil {
				b.FailNow()
			}
		}
	}
}

func cleanUpAfterBenchmark(
	b *testing.B,
	eventstore *eventstore.CustomerEventStore,
	id values.CustomerID,
) {

	if err := eventstore.Delete(id); err != nil {
		b.FailNow()
	}
}
