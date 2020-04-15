package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeName(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
		personName := value.RebuildPersonName("Kevin", "Ball")
		changedPersonName := value.RebuildPersonName("Latoya", "Ball")

		customerWasRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		changeName, err := domain.BuildChangeCustomerName(
			customerID.String(),
			changedPersonName.GivenName(),
			changedPersonName.FamilyName(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's name", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					recordedEvents, err = customer.ChangeName(eventStream, changeName)
					So(err, ShouldBeNil)

					Convey("Then CustomerNameChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						nameChanged, ok := recordedEvents[0].(domain.CustomerNameChanged)
						So(ok, ShouldBeTrue)
						So(nameChanged, ShouldNotBeNil)
						So(nameChanged.CustomerID().Equals(customerID), ShouldBeTrue)
						So(nameChanged.PersonName().Equals(changedPersonName), ShouldBeTrue)
						So(nameChanged.IsFailureEvent(), ShouldBeFalse)
						So(nameChanged.FailureReason(), ShouldBeNil)
						So(nameChanged.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's name to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					changeName, err = domain.BuildChangeCustomerName(
						customerID.String(),
						personName.GivenName(),
						personName.FamilyName(),
					)
					So(err, ShouldBeNil)

					recordedEvents, err = customer.ChangeName(eventStream, changeName)
					So(err, ShouldBeNil)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's name to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("and CustomerNameChanged", func() {
					nameChanged := domain.BuildCustomerNameChanged(
						customerID,
						changedPersonName,
						2,
					)

					eventStream = append(eventStream, nameChanged)

					Convey("When ChangeCustomerName", func() {
						recordedEvents, err = customer.ChangeName(eventStream, changeName)
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
				eventStream := es.EventStream{customerWasRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						domain.BuildCustomerDeleted(customerID, emailAddress, 2),
					)

					Convey("When ChangeCustomerName", func() {
						_, err := customer.ChangeName(eventStream, changeName)

						Convey("Then it should report an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
