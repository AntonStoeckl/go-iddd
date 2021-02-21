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

func TestChangeName(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := value.GenerateCustomerID()
		emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		personName, err := value.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)
		changedPersonName, err := value.BuildPersonName("Latoya", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildChangeCustomerName(customerID, changedPersonName)
		commandWithOriginalName := domain.BuildChangeCustomerName(customerID, personName)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			personName,
			es.GenerateMessageID(),
			1,
		)

		customerDeleted := domain.BuildCustomerDeleted(
			customerID,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO 1: Change a Customer's name", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerName", func() {
					recordedEvents, err = customer.ChangeName(eventStream, command)
					So(err, ShouldBeNil)

					Convey("Then CustomerNameChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerNameChanged)
						So(ok, ShouldBeTrue)
						So(event, ShouldNotBeNil)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.PersonName().Equals(changedPersonName), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeFalse)
						So(event.FailureReason(), ShouldBeNil)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's name to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerName", func() {
					recordedEvents, err = customer.ChangeName(eventStream, commandWithOriginalName)
					So(err, ShouldBeNil)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's name to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerNameChanged", func() {
					nameChanged := domain.BuildCustomerNameChanged(
						customerID,
						changedPersonName,
						es.GenerateMessageID(),
						2,
					)

					eventStream = append(eventStream, nameChanged)

					Convey("When ChangeCustomerName", func() {
						recordedEvents, err = customer.ChangeName(eventStream, command)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Try to change a Customer's name when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						customerDeleted,
					)

					Convey("When ChangeCustomerName", func() {
						_, err := customer.ChangeName(eventStream, command)

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
