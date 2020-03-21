package application_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
)

func BenchmarkCustomerScenario(b *testing.B) {
	var err error

	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()

	emailAddress := "fiona@gallagher.net"
	givenName := "Fiona"
	familyName := "Galagher"
	newEmailAddress := "fiona@pratt.net"
	newGivenName := "Fiona"
	newFamilyName := "Pratt"

	registerCustomer, err := commands.BuildRegisterCustomer(
		emailAddress,
		givenName,
		familyName,
	)

	if err != nil {
		b.FailNow()
	}

	confirmCustomerEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
		registerCustomer.CustomerID().ID(),
		registerCustomer.ConfirmationHash().Hash(),
	)

	if err != nil {
		b.FailNow()
	}

	changeCustomerEmailAddress, err := commands.BuildChangeCustomerEmailAddress(
		registerCustomer.CustomerID().ID(),
		newEmailAddress,
	)

	if err != nil {
		b.FailNow()
	}

	changeCustomerEmailAddressBack, err := commands.BuildChangeCustomerEmailAddress(
		registerCustomer.CustomerID().ID(),
		emailAddress,
	)

	if err != nil {
		b.FailNow()
	}

	changeCustomerName, err := commands.BuildChangeCustomerName(
		registerCustomer.CustomerID().ID(),
		newGivenName,
		newFamilyName,
	)

	if err != nil {
		b.FailNow()
	}

	changeCustomerNameBack, err := commands.BuildChangeCustomerName(
		registerCustomer.CustomerID().ID(),
		givenName,
		familyName,
	)

	if err != nil {
		b.FailNow()
	}

	prepareForBenchmark(b,
		commandHandler,
		registerCustomer,
		confirmCustomerEmailAddress,
		changeCustomerEmailAddress,
		changeCustomerEmailAddressBack,
	)

	/**************************************************/

	b.Run("ChangeName", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if n%2 == 0 {
				if err = commandHandler.ChangeCustomerName(changeCustomerName); err != nil {
					b.FailNow()
				}
			} else {
				if err = commandHandler.ChangeCustomerName(changeCustomerNameBack); err != nil {
					b.FailNow()
				}
			}
		}
	})

	/**************************************************/

	b.Run("CustomerViewByID", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err = queryHandler.CustomerViewByID(registerCustomer.CustomerID()); err != nil {
				b.FailNow()
			}
		}
	})

	cleanUpAfterBenchmark(
		b,
		diContainer.GetCustomerEventStore(),
		registerCustomer.CustomerID(),
	)
}

func prepareForBenchmark(
	b *testing.B,
	commandHandler *command.CustomerCommandHandler,
	registerCustomer commands.RegisterCustomer,
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress,
	changeCustomerEmailAddress commands.ChangeCustomerEmailAddress,
	changeCustomerEmailAddressBack commands.ChangeCustomerEmailAddress,
) {

	var err error

	if err = commandHandler.RegisterCustomer(registerCustomer); err != nil {
		b.FailNow()
	}

	if err = commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress); err != nil {
		b.FailNow()
	}

	for n := 0; n < 100; n++ {
		if n%2 == 0 {
			if err = commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress); err != nil {
				b.FailNow()
			}
		} else {
			if err = commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddressBack); err != nil {
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
