package readmodel_test

import (
	"database/sql"
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/writemodel"
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
	"go-iddd/service/customer/application/writemodel/domain/customer/events"
	"go-iddd/service/lib"
	"go-iddd/service/lib/eventstore/postgres/database"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerQueryHandlerScenarios(t *testing.T) {
	diContainer := setUpDiContainerForCustomerQueryHandlerScenarios()
	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerView customer.View

		emailAddress := "fiona@gallagher.net"
		givenName := "Fiona"
		familyName := "Galagher"
		newEmailAddress := "fiona@pratt.net"

		registerCustomer, err := commands.BuildRegisterCustomer(
			emailAddress,
			givenName,
			familyName,
		)
		So(err, ShouldBeNil)

		customerID := customer.RebuildID(registerCustomer.CustomerID().ID())

		confirmCustomerEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			registerCustomer.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		changeCustomerEmailAddress, err := commands.BuildChangeCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			newEmailAddress,
		)
		So(err, ShouldBeNil)

		expectedCustomerView := customer.View{
			ID:                      customerID.ID(),
			EmailAddress:            emailAddress,
			IsEmailAddressConfirmed: false,
			GivenName:               givenName,
			FamilyName:              familyName,
			Version:                 1,
		}

		Convey("\nSCENARIO 1: Retrieving a registered Customer", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(registerCustomer, commandHandler)

				Convey("When the Customer is retrieved by ID", func() {
					customerView, err = queryHandler.CustomerViewByID(customerID)
					So(err, ShouldBeNil)

					Convey("Then the Customer view should be as expected", func() {
						So(customerView, ShouldResemble, expectedCustomerView)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Retrieving a registered Customer with a confirmed emailAddress", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey("When the Customer is retrieved by ID", func() {
						customerView, err = queryHandler.CustomerViewByID(customerID)
						So(err, ShouldBeNil)

						Convey("Then the Customer view should be as expected", func() {
							expectedCustomerView.IsEmailAddressConfirmed = true
							expectedCustomerView.Version = 2

							So(customerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Retrieving a registered Customer with an emailAddress that was confirmed and then changed", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey("And given she changed her emailAddress", func() {
						GivenEmailAddressChanged(changeCustomerEmailAddress, commandHandler)

						Convey("When the Customer is retrieved by ID", func() {
							customerView, err = queryHandler.CustomerViewByID(customerID)
							So(err, ShouldBeNil)

							Convey("Then the Customer view should be as expected", func() {
								expectedCustomerView.EmailAddress = newEmailAddress
								expectedCustomerView.Version = 3

								So(customerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Trying to retrieve a Customer which does not exist", func() {
			Convey("When the Customer is retrieved by ID", func() {
				customerView, err := queryHandler.CustomerViewByID(customer.GenerateID())

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(customerView, ShouldBeZeroValue)
				})
			})
		})

		Reset(func() {
			err := diContainer.GetCustomerEventStoreForWriteModel().Delete(registerCustomer.CustomerID())
			So(err, ShouldBeNil)
		})
	})
}

func GivenCustomerRegistered(
	registerCustomer commands.RegisterCustomer,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.RegisterCustomer(registerCustomer)
	So(err, ShouldBeNil)

}

func GivenEmailAddressConfirmed(
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func GivenEmailAddressChanged(
	changeCustomerEmailAddress commands.ChangeCustomerEmailAddress,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func setUpDiContainerForCustomerQueryHandlerScenarios() *cmd.DIContainer {
	config, err := cmd.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		panic(err)
	}

	err = migrator.Up()
	if err != nil {
		panic(err)
	}

	diContainer, err := cmd.NewDIContainer(
		db,
		events.UnmarshalCustomerEvent,
		customer.UnmarshalCustomerEvent,
	)

	if err != nil {
		panic(err)
	}

	return diContainer
}
