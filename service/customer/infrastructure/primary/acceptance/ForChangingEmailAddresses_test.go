package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_ForChangingEmailAddresses(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)

		commandHandler := application.NewCommandHandler(customerEventStore)

		registered := events.ItWasRegistered(
			values.GenerateCustomerID(),
			values.RebuildEmailAddress("john@doe.com"),
			values.RebuildConfirmationHash("john@doe.com"),
			values.RebuildPersonName("John", "Doe"),
			uint(1),
		)

		emailAddressConfirmed := events.EmailAddressWasChanged(
			registered.CustomerID(),
			registered.EmailAddress(),
			registered.ConfirmationHash(),
			uint(2),
		)

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			registered.CustomerID().ID(),
			"john+change@doe.com",
		)
		So(err, ShouldBeNil)

		emailAddressChanged := events.EmailAddressWasChanged(
			changeEmailAddress.CustomerID(),
			changeEmailAddress.EmailAddress(),
			changeEmailAddress.ConfirmationHash(),
			uint(3),
		)

		containsOnlyEmailAddressChangedEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.EmailAddressChanged)

			return ok
		}

		containsOnlyEmailAddressConfirmedEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.EmailAddressConfirmed)

			return ok
		}

		Convey("Given a registered Customer", func() {
			Convey("And given their emailAddress has been confirmed", func() {
				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered, emailAddressConfirmed}, nil).
					Once()

				Convey("When the emailAddress is changed", func() {
					customerEventStore.
						On(
							"Add",
							mock.MatchedBy(containsOnlyEmailAddressChangedEvent),
							registered.CustomerID(),
						).
						Return(nil).
						Once()

					err = commandHandler.ChangeEmailAddress(changeEmailAddress)

					Convey("It should succeed", func() {
						So(err, ShouldBeNil)

						Convey("And the changed emailAddress should be unconfirmed", func() {
							customerEventStore.
								On("EventStreamFor", registered.CustomerID()).
								Return(es.DomainEvents{registered, emailAddressConfirmed, emailAddressChanged}, nil).
								Once()

							customerEventStore.
								On(
									"Add",
									mock.MatchedBy(containsOnlyEmailAddressConfirmedEvent),
									registered.CustomerID(),
								).
								Return(nil).
								Once()

							confirmEmailAddress, err := commands.NewConfirmEmailAddress(
								changeEmailAddress.CustomerID().ID(),
								changeEmailAddress.EmailAddress().EmailAddress(),
								changeEmailAddress.ConfirmationHash().Hash(),
							)
							So(err, ShouldBeNil)

							err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
							So(err, ShouldBeNil)
						})
					})
				})
			})

			Convey("And given their emailAddress has already been changed", func() {
				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered, emailAddressChanged}, nil).
					Once()

				Convey("When the emailAddress is changed to the same value again", func() {
					customerEventStore.
						On(
							"Add",
							es.DomainEvents(nil),
							registered.CustomerID(),
						).
						Return(nil).
						Once()

					err = commandHandler.ChangeEmailAddress(changeEmailAddress)

					Convey("It should be ignored", func() {
						So(err, ShouldBeNil)
					})
				})
			})

			Convey("When the emailAddress is changed with an invalid command", func() {
				err := commandHandler.ChangeEmailAddress(commands.ChangeEmailAddress{})

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			Convey("Assuming that the recordedEvents can't be added", func() {
				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered}, nil).
					Once()

				customerEventStore.
					On("Add", mock.Anything, mock.Anything, mock.Anything).
					Return(lib.ErrTechnical).
					Once()

				Convey("When the emailAddress is confirmed", func() {
					err := commandHandler.ChangeEmailAddress(changeEmailAddress)

					Convey("It should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					})
				})
			})
		})

		Convey("Given an unregistered Customer", func() {
			customerID := values.GenerateCustomerID()

			customerEventStore.
				On("EventStreamFor", customerID).
				Return(es.DomainEvents{}, lib.ErrNotFound).
				Once()

			Convey("When the emailAddress is confirmed", func() {
				changeEmailAddress, err := commands.NewChangeEmailAddress(
					customerID.ID(),
					"john@doe.com",
				)
				So(err, ShouldBeNil)

				err = commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}
