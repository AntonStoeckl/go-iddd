package application_test

import (
	"fmt"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/eventstore"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommandHandlerScenarios(t *testing.T) {
	Convey("Customer Lifecycle Scenarios", t, func() {
		diContainer, err := infrastructure.SetUpDIContainer()
		So(err, ShouldBeNil)

		commandHandler := diContainer.GetCustomerCommandHandler()
		customerEventStore := diContainer.GetCustomerEventStore()

		emailAddress := "fiona@gallagher.net"
		givenName := "Fiona"
		familyName := "Gallagher"

		register, err := commands.NewRegister(
			emailAddress,
			givenName,
			familyName,
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: A prospective Customer registers", func() {
			Convey(fmt.Sprintf("When I register as %s %s with %s", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)

				Convey("Then I should have an account with an unconfirmed email address", func() {
					So(err, ShouldBeNil)
					MyAccountShouldBeRegistered(register.CustomerID(), customerEventStore)
				})
			})
		})

		Convey("\nSCENARIO 2: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given I registered as %s %s with %s", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("When I confirm my email address with a valid confirmation hash", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						register.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then my email address should be confirmed", func() {
						MyEmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)

						Convey("When I try to confirm it again", func() {
							err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

							Convey("Then it should be ignored", func() {
								So(err, ShouldBeNil)
								MyEmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: A Customer fails to confirm her email address", func() {
			Convey(fmt.Sprintf("Given I registered as %s %s with %s", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("When I try to confirm my email address with an invalid confirmation hash", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						values.GenerateConfirmationHash(register.EmailAddress().EmailAddress()).Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

					Convey("Then it should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("and my email address should be unconfirmed", func() {
							MyEmailAddressShouldNotBeConfirmed(register.CustomerID(), customerEventStore)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: A Customer changes her confirmed email address", func() {
			Convey(fmt.Sprintf("Given I registered as %s %s with %s", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("and I confirmed my email address", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						register.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
					So(err, ShouldBeNil)

					Convey("When I change my email address", func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							"john@doe.com",
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then my email address should be changed and unconfirmed", func() {
							MyEmailAddressShouldBeChangedAndUnconfirmed(register.CustomerID(), customerEventStore)

							Convey("When I try to change it again", func() {
								changeEmailAddress, err := commands.NewChangeEmailAddress(
									register.CustomerID().ID(),
									"john@doe.com",
								)
								So(err, ShouldBeNil)

								err = commandHandler.ChangeEmailAddress(changeEmailAddress)
								Convey("Then it should be ignored", func() {
									So(err, ShouldBeNil)
									MyEmailAddressShouldBeChangedAndUnconfirmed(register.CustomerID(), customerEventStore)
								})
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: A Customer confirms her changed email address", func() {
			Convey(fmt.Sprintf("Given I registered as %s %s with %s", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("and I confirmed my email address", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						register.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
					So(err, ShouldBeNil)

					Convey("and I changed my email address", func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							"john@doe.com",
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("When I confirm my changed email address", func() {
							confirmEmailAddress, err := commands.NewConfirmEmailAddress(
								changeEmailAddress.CustomerID().ID(),
								changeEmailAddress.EmailAddress().EmailAddress(),
								changeEmailAddress.ConfirmationHash().Hash(),
							)
							So(err, ShouldBeNil)

							err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then my changed email address should be confirmed", func() {
								MyChangedEmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)
							})
						})
					})
				})
			})

		})

		Reset(func() {
			err := customerEventStore.Delete(register.CustomerID())
			So(err, ShouldBeNil)
		})
	})
}

func MyAccountShouldBeRegistered(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 1)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
}

func MyEmailAddressShouldBeConfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 2)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
	So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressConfirmed{})
}

func MyEmailAddressShouldNotBeConfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 2)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
	So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressConfirmationFailed{})
}

func MyChangedEmailAddressShouldBeConfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 4)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
	So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressConfirmed{})
	So(eventStream[2], ShouldHaveSameTypeAs, events.EmailAddressChanged{})
	So(eventStream[3], ShouldHaveSameTypeAs, events.EmailAddressConfirmed{})
}

func MyEmailAddressShouldBeChangedAndUnconfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 3)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
	So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressConfirmed{})
	So(eventStream[2], ShouldHaveSameTypeAs, events.EmailAddressChanged{})
}
