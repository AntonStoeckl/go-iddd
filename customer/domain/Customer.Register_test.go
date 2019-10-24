package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterCustomer(t *testing.T) {
	Convey("When a Customer is registered", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		register, err := commands.NewRegister(
			id.String(),
			emailAddress.EmailAddress(),
			personName.GivenName(),
			personName.FamilyName(),
		)
		So(err, ShouldBeNil)

		recordedEvents := domain.RegisterCustomer(register)

		Convey("It should record CustomerRegistered", func() {
			So(recordedEvents, ShouldHaveLength, 1)
			registered, ok := recordedEvents[0].(*events.Registered)
			So(ok, ShouldBeTrue)
			So(registered, ShouldNotBeNil)
			So(registered.ID().Equals(id), ShouldBeTrue)
			So(registered.ConfirmableEmailAddress().Equals(emailAddress), ShouldBeTrue)
			_, err := values.RebuildConfirmationHash(registered.ConfirmableEmailAddress().ConfirmationHash())
			So(err, ShouldBeNil)
			So(registered.PersonName().Equals(personName), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})

		Convey("And it should reconstitute a Customer", func() {
			customer, err := domain.ReconstituteCustomerFrom(recordedEvents)
			So(err, ShouldBeNil)
			So(customer, ShouldNotBeNil)
			So(customer, ShouldHaveSameTypeAs, (*domain.Customer)(nil))
			So(customer.ID(), ShouldImplement, (*shared.IdentifiesAggregates)(nil))
			So(customer.ID().Equals(id), ShouldBeTrue)
		})
	})
}
