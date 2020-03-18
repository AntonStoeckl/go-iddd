package application_test

import (
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/readmodel"
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/writemodel"
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerScenarios(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var actualCustomerView customer.View

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

		customerID := registerCustomer.CustomerID()
		customerViewID := customer.RebuildID(registerCustomer.CustomerID().ID())

		confirmCustomerEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			registerCustomer.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		confirmCustomerEmailAddressWithInvalidHash, err := commands.BuildConfirmCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			values.GenerateConfirmationHash(emailAddress).Hash(),
		)
		So(err, ShouldBeNil)

		changeCustomerEmailAddress, err := commands.BuildChangeCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			newEmailAddress,
		)
		So(err, ShouldBeNil)

		confirmChangedCustomerEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			registerCustomer.CustomerID().ID(),
			changeCustomerEmailAddress.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		expectedCustomerView := customer.View{
			ID:                      customerViewID.ID(),
			EmailAddress:            emailAddress,
			IsEmailAddressConfirmed: false,
			GivenName:               givenName,
			FamilyName:              familyName,
			Version:                 1,
		}

		Convey("\nSCENARIO: A prospective Customer registers her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("When she retrieves here account data", func() {
					actualCustomerView = retrieveAccountData(queryHandler, customerViewID)

					Convey("Then she should see an unconfirmed account with the data she supplied", func() {
						So(actualCustomerView, ShouldResemble, expectedCustomerView)
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey("When she retrieves here account data", func() {
						actualCustomerView = retrieveAccountData(queryHandler, customerViewID)

						Convey("Then she should see that here account is now confirmed", func() {
							expectedCustomerView.IsEmailAddressConfirmed = true
							expectedCustomerView.Version = 2

							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer fails to confirm her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("When she tries to confirm her email address with a wrong confirmation hash", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddressWithInvalidHash)

					Convey("Then she should receive an error", func() {
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("And she should see that her account is still unconfirmed", func() {
							actualCustomerView, err = queryHandler.CustomerViewByID(customerViewID)
							So(err, ShouldBeNil)
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer changes her (confirmed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey(fmt.Sprintf("And given she changed her email address to [%s]", newEmailAddress), func() {
						givenCustomerEmailAddressWasChanged(changeCustomerEmailAddress, commandHandler)

						Convey("When she retrieves here account data", func() {
							actualCustomerView = retrieveAccountData(queryHandler, customerViewID)

							Convey("Then she should see that her email address is changed and unconfirmed", func() {
								expectedCustomerView.EmailAddress = newEmailAddress
								expectedCustomerView.Version = 3

								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her (changed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey(fmt.Sprintf("And given she changed her email address to [%s]", newEmailAddress), func() {
						givenCustomerEmailAddressWasChanged(changeCustomerEmailAddress, commandHandler)

						Convey("And given she confirmed her changed email address", func() {
							givenCustomerEmailAddressWasConfirmed(confirmChangedCustomerEmailAddress, commandHandler)

							Convey("When she retrieves here account data", func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerViewID)

								Convey("Then she should see that her changed email address is confirmed", func() {
									expectedCustomerView.EmailAddress = newEmailAddress
									expectedCustomerView.IsEmailAddressConfirmed = true
									expectedCustomerView.Version = 4

									So(actualCustomerView, ShouldResemble, expectedCustomerView)
								})
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer deletes her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she deleted here account", func() {
					givenCustomerAccountWasDeleted(diContainer, customerID)

					Convey("When she tries to retrieve here account data", func() {
						customerView, err := queryHandler.CustomerViewByID(customer.GenerateID())

						Convey("Then it should fail", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(customerView, ShouldBeZeroValue)
						})
					})
				})
			})
		})

		Reset(func() {
			err := diContainer.GetCustomerEventStoreForWriteModel().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func givenCustomerRegistered(
	registerCustomer commands.RegisterCustomer,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.RegisterCustomer(registerCustomer)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasConfirmed(
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	changeCustomerEmailAddress commands.ChangeCustomerEmailAddress,
	commandHandler *writemodel.CustomerCommandHandler,
) {

	err := commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func givenCustomerAccountWasDeleted(
	diContainer *cmd.DIContainer,
	customerID values.CustomerID,
) {

	// TODO: introduce a command to (soft?) delete an account
	err := diContainer.GetCustomerEventStoreForWriteModel().Delete(customerID)
	So(err, ShouldBeNil)
}

func retrieveAccountData(
	queryHandler *readmodel.CustomerQueryHandler,
	id customer.ID,
) customer.View {

	customerView, err := queryHandler.CustomerViewByID(id)
	So(err, ShouldBeNil)

	return customerView
}
