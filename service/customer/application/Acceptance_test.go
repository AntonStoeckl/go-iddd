package application_test

import (
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/command"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/customer"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/application/query"
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
		newGivenName := "Fiona"
		newFamilyName := "Pratt"

		registerCustomer, err := commands.BuildRegisterCustomer(
			emailAddress,
			givenName,
			familyName,
		)
		So(err, ShouldBeNil)

		customerID := registerCustomer.CustomerID()

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

		changeCustomerName, err := commands.BuildChangeCustomerName(
			registerCustomer.CustomerID().ID(),
			newGivenName,
			newFamilyName,
		)

		expectedCustomerView := customer.View{
			ID:                      customerID.ID(),
			EmailAddress:            emailAddress,
			IsEmailAddressConfirmed: false,
			GivenName:               givenName,
			FamilyName:              familyName,
			Version:                 1,
		}

		Convey("\nSCENARIO: A prospective Customer registers her account", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err = commandHandler.RegisterCustomer(registerCustomer)
				So(err, ShouldBeNil)

				details := fmt.Sprintf("\n\tGivenName: %s", givenName)
				details += fmt.Sprintf("\n\tFamilyName: %s", familyName)
				details += fmt.Sprintf("\n\tEmailAddress: %s (unconfirmed)", emailAddress)

				Convey(fmt.Sprintf("Then her account should show the data she supplied: %s", details), func() {
					actualCustomerView = retrieveAccountData(queryHandler, customerID)
					So(actualCustomerView, ShouldResemble, expectedCustomerView)
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("When she confirms her email address", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then her email address should be confirmed", func() {
						actualCustomerView = retrieveAccountData(queryHandler, customerID)
						expectedCustomerView.IsEmailAddressConfirmed = true
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey("And when she confirms her email address again", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then her email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								expectedCustomerView.IsEmailAddressConfirmed = true
								expectedCustomerView.Version = 2
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
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

						Convey("And her email address should still be unconfirmed", func() {
							actualCustomerView = retrieveAccountData(queryHandler, customerID)
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer fails to confirm her already confirmed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(confirmCustomerEmailAddress, commandHandler)

					Convey("When she tries to confirm her email address again with a wrong confirmation hash", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddressWithInvalidHash)

						Convey("Then she should receive an error", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("And her email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								expectedCustomerView.IsEmailAddressConfirmed = true
								expectedCustomerView.Version = 3
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
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

					Convey(fmt.Sprintf("When she changes her email address to [%s]", newEmailAddress), func() {
						err = commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("Then her email address should be [%s] and unconfirmed", newEmailAddress), func() {
							actualCustomerView = retrieveAccountData(queryHandler, customerID)
							expectedCustomerView.EmailAddress = newEmailAddress
							expectedCustomerView.Version = 3
							So(actualCustomerView, ShouldResemble, expectedCustomerView)

							Convey(fmt.Sprintf("And when she tries to change her email address to [%s] again", newEmailAddress), func() {
								err = commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
								So(err, ShouldBeNil)

								Convey(fmt.Sprintf("Then her email address should still be [%s]", newEmailAddress), func() {
									actualCustomerView = retrieveAccountData(queryHandler, customerID)
									expectedCustomerView.EmailAddress = newEmailAddress
									expectedCustomerView.Version = 3
									So(actualCustomerView, ShouldResemble, expectedCustomerView)
								})
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

						Convey("When she confirms her changed email address", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(confirmChangedCustomerEmailAddress)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her email address should be [%s] and confirmed", newEmailAddress), func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
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

		Convey("\nSCENARIO: A Customer changes her name", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey(fmt.Sprintf("When she changes her name to [%s %s]", newGivenName, newFamilyName), func() {
					err = commandHandler.ChangeCustomerName(changeCustomerName)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then her name should be [%s %s]", newGivenName, newFamilyName), func() {
						actualCustomerView = retrieveAccountData(queryHandler, customerID)
						expectedCustomerView.GivenName = newGivenName
						expectedCustomerView.FamilyName = newFamilyName
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey(fmt.Sprintf("And when she tries to change her name to [%s %s] again", newGivenName, newFamilyName), func() {
							err = commandHandler.ChangeCustomerName(changeCustomerName)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her name should still be [%s %s]", newGivenName, newFamilyName), func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								expectedCustomerView.GivenName = newGivenName
								expectedCustomerView.FamilyName = newFamilyName
								expectedCustomerView.Version = 2
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer deletes her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				givenCustomerRegistered(registerCustomer, commandHandler)

				Convey("When she deletes her account", func() {
					// TODO: introduce a command to (soft?) delete an account
					err = diContainer.GetCustomerEventStore().Delete(customerID)
					So(err, ShouldBeNil)

					Convey("And when she tries to retrieve her account data", func() {
						customerView, err := queryHandler.CustomerViewByID(customerID)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(customerView, ShouldBeZeroValue)
						})
					})

					Convey("And when she tries to confirm her email address", func() {
						err := commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her email address", func() {
						err := commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her name", func() {
						err := commandHandler.ChangeCustomerName(changeCustomerName)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err := diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func givenCustomerRegistered(
	registerCustomer commands.RegisterCustomer,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.RegisterCustomer(registerCustomer)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasConfirmed(
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	changeCustomerEmailAddress commands.ChangeCustomerEmailAddress,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func retrieveAccountData(
	queryHandler *query.CustomerQueryHandler,
	id values.CustomerID,
) customer.View {

	customerView, err := queryHandler.CustomerViewByID(id)
	So(err, ShouldBeNil)

	return customerView
}
