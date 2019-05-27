package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestRegisterCustomer(t *testing.T) {
	Convey("When a Customer is registered", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.Register(register)

		Convey("It should succeed", func() {
			So(customer, ShouldNotBeNil)
			So(customer, ShouldImplement, (*domain.Customer)(nil))

			Convey("And it should expose the expected AggregateIdentifier", func() {
				So(customer.AggregateIdentifier().String(), ShouldEqual, id)
			})

			Convey("And it should expose the expected AggregateName", func() {
				So(customer.AggregateName(), ShouldEqual, "Customer")
			})

			Convey("And it should record that a Customer was registered", func() {
				recordedEvents := customer.RecordedEvents()
				registered := findCustomerEventIn(
					recordedEvents,
					new(events.Registered),
				).(*events.Registered)

				So(registered, ShouldNotBeNil)
				So(registered.ID().String(), ShouldEqual, id)
				So(registered.ConfirmableEmailAddress().EmailAddress(), ShouldEqual, emailAddress)
				So(registered.ConfirmableEmailAddress().ConfirmationHash(), ShouldNotBeBlank)
				So(registered.PersonName().GivenName(), ShouldEqual, givenName)
				So(registered.PersonName().FamilyName(), ShouldEqual, familyName)

				Convey("And it should not record anything else", func() {
					So(recordedEvents, ShouldHaveLength, 1)
				})
			})

			Convey("And it should not apply further Register commands", func() {
				err = customer.Apply(register)
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandCanNotBeHandled), ShouldBeTrue)
			})
		})
	})
}
