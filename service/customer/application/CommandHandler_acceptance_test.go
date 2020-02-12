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

		issuedHash := register.ConfirmationHash().Hash()

		Convey("\nSCENARIO 1: A prospective Customer registers", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)

				Convey("Then she should have an unconfirmed account", func() {
					So(err, ShouldBeNil)
					AccountShouldBeRegisteredAndUnconfirmed(register.CustomerID(), customerEventStore)
				})
			})
		})

		Convey("\nSCENARIO 2: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", issuedHash), func() {
					Convey(fmt.Sprintf("When she confirms her email address with confirmation hash [%s]", issuedHash), func() {
						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							register.CustomerID().ID(),
							register.EmailAddress().EmailAddress(),
							issuedHash,
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then her email address should be confirmed", func() {
							EmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)

							Convey(fmt.Sprintf("When she tries to confirm it again with confirmation hash [%s]", issuedHash), func() {
								err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

								Convey("Then it should be ignored", func() {
									So(err, ShouldBeNil)
									EmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)
								})
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: A Customer fails to confirm her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", issuedHash), func() {
					invalidHash := values.GenerateConfirmationHash(register.EmailAddress().EmailAddress()).Hash()

					Convey(fmt.Sprintf("When she tries to confirm her email address with invalid confirmation hash [%s]", invalidHash), func() {
						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							register.CustomerID().ID(),
							register.EmailAddress().EmailAddress(),
							invalidHash,
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

						Convey("Then it should fail", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("and her email address should be unconfirmed", func() {
								EmailAddressShouldNotBeConfirmed(register.CustomerID(), customerEventStore)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: A Customer changes her confirmed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("and she confirmed her email address", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						register.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
					So(err, ShouldBeNil)

					changedEmailAddress := "fiona@pratt.net"

					Convey(fmt.Sprintf("When she changes her email address to [%s]", changedEmailAddress), func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							changedEmailAddress,
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then her email address should be changed and unconfirmed", func() {
							MyEmailAddressShouldBeChangedAndUnconfirmed(register.CustomerID(), customerEventStore)

							Convey(fmt.Sprintf("When she tries to change it again to [%s]", changedEmailAddress), func() {
								changeEmailAddress, err := commands.NewChangeEmailAddress(
									register.CustomerID().ID(),
									changedEmailAddress,
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
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("and she confirmed her email address", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						register.CustomerID().ID(),
						register.EmailAddress().EmailAddress(),
						register.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
					So(err, ShouldBeNil)

					changedEmailAddress := "fiona@pratt.net"

					Convey(fmt.Sprintf("and she changed her email address to [%s]", changedEmailAddress), func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							changedEmailAddress,
						)
						So(err, ShouldBeNil)

						issuedHash := changeEmailAddress.ConfirmationHash().Hash()

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", issuedHash), func() {
							Convey(fmt.Sprintf("When she confirms her changed email address with confirmation hash [%s]", issuedHash), func() {
								confirmEmailAddress, err := commands.NewConfirmEmailAddress(
									changeEmailAddress.CustomerID().ID(),
									changeEmailAddress.EmailAddress().EmailAddress(),
									issuedHash,
								)
								So(err, ShouldBeNil)

								err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
								So(err, ShouldBeNil)

								Convey("Then her changed email address should be confirmed", func() {
									MyChangedEmailAddressShouldBeConfirmed(register.CustomerID(), customerEventStore)
								})
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

func AccountShouldBeRegisteredAndUnconfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 1)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
}

func EmailAddressShouldBeConfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)
	So(eventStream, ShouldHaveLength, 2)
	So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
	So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressConfirmed{})
}

func EmailAddressShouldNotBeConfirmed(customerID values.CustomerID, customerEventStore *eventstore.CustomerEventStore) {
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
