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

func Test_ForConfirmingEmailAddresses(t *testing.T) {
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

		emailAddressConfirmed := events.EmailAddressWasConfirmed(
			registered.CustomerID(),
			registered.EmailAddress(),
			uint(2),
		)

		withValidHash, err := commands.NewConfirmEmailAddress(
			registered.CustomerID().ID(),
			registered.EmailAddress().EmailAddress(),
			registered.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		withInvalidHash, err := commands.NewConfirmEmailAddress(
			registered.CustomerID().ID(),
			registered.EmailAddress().EmailAddress(),
			"some_invalid_hash",
		)
		So(err, ShouldBeNil)

		containsOnlyEmailAddressConfirmedEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.EmailAddressConfirmed)

			return ok
		}

		containsOnlyEmailAddressConfirmationFailedEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.EmailAddressConfirmationFailed)

			return ok
		}

		Convey("Given a registered Customer", func() {
			Convey("And given their emailAddress has not been confirmed", func() {
				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered}, nil).
					Once()

				Convey("When the emailAddress is confirmed with a valid confirmationHash", func() {
					customerEventStore.
						On(
							"Add",
							mock.MatchedBy(containsOnlyEmailAddressConfirmedEvent),
							registered.CustomerID(),
						).
						Return(nil).
						Once()

					err = commandHandler.ConfirmEmailAddress(withValidHash)

					Convey("It should succeed", func() {
						So(err, ShouldBeNil)
					})
				})

				Convey("When the emailAddress is confirmed with an invalid confirmationHash", func() {
					customerEventStore.
						On(
							"Add",
							mock.MatchedBy(containsOnlyEmailAddressConfirmationFailedEvent),
							registered.CustomerID(),
						).
						Return(nil).
						Once()

					err = commandHandler.ConfirmEmailAddress(withInvalidHash)

					Convey("It should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)
					})
				})
			})

			Convey("And given their emailAddress has already been confirmed", func() {
				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered, emailAddressConfirmed}, nil).
					Once()

				Convey("When the emailAddress is confirmed again", func() {
					customerEventStore.
						On(
							"Add",
							es.DomainEvents(nil),
							registered.CustomerID(),
						).
						Return(nil).
						Once()

					err = commandHandler.ConfirmEmailAddress(withValidHash)

					Convey("It should succeed", func() {
						So(err, ShouldBeNil)
					})
				})
			})

			Convey("When the emailAddress is confirmed with an invalid command", func() {
				err := commandHandler.ConfirmEmailAddress(commands.ConfirmEmailAddress{})

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
					err := commandHandler.ConfirmEmailAddress(withValidHash)

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
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					customerID.ID(),
					"john@doe.com",
					"any_hash",
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}
