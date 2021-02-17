package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
		personName := value.RebuildPersonName("Kevin", "Ball")
		changedEmailAddress := value.RebuildEmailAddress("latoya@ball.net")

		command, err := domain.BuildChangeCustomerEmailAddress(
			customerID.String(),
			changedEmailAddress.String(),
		)
		So(err, ShouldBeNil)

		commandWithOriginalEmailAddress, err := domain.BuildChangeCustomerEmailAddress(
			customerID.String(),
			emailAddress.String(),
		)
		So(err, ShouldBeNil)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			es.GenerateMessageID(),
			1,
		)

		customerEmailAddressChanged := domain.BuildCustomerEmailAddressChanged(
			customerID,
			changedEmailAddress,
			command.ConfirmationHash(),
			emailAddress,
			es.GenerateMessageID(),
			2,
		)

		customerDeleted := domain.BuildCustomerDeleted(
			customerID,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO 1: Change a Customer's emailAddress", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerEmailAddress", func() {
					recordedEvents, err = customer.ChangeEmailAddress(eventStream, command)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerEmailAddressChanged)
						So(ok, ShouldBeTrue)
						So(event, ShouldNotBeNil)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.EmailAddress().Equals(changedEmailAddress), ShouldBeTrue)
						So(event.ConfirmationHash().Equals(command.ConfirmationHash()), ShouldBeTrue)
						So(event.PreviousEmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeFalse)
						So(event.FailureReason(), ShouldBeNil)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's emailAddress to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerEmailAddress", func() {
					recordedEvents, err = customer.ChangeEmailAddress(eventStream, commandWithOriginalEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's emailAddress to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressChanged", func() {
					eventStream = append(eventStream, customerEmailAddressChanged)

					Convey("When ChangeCustomerEmailAddress", func() {
						recordedEvents, err = customer.ChangeEmailAddress(eventStream, command)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Try to change a Customer's emailAddress when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(eventStream, customerDeleted)

					Convey("When ChangeCustomerEmailAddress", func() {
						_, err := customer.ChangeEmailAddress(eventStream, command)

						Convey("Then it should report an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
