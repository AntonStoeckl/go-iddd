package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artefacts", t, func() {
		register, err := commands.NewRegister(
			"kevin@ball.com",
			"Kevin",
			"Ball",
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				recordedEvents := domain.RegisterCustomer(register)

				Convey("Then CustomerRegistered", func() {
					ThenCustomerRegistered(recordedEvents, register)
				})
			})
		})
	})
}

func ThenCustomerRegistered(recordedEvents es.DomainEvents, register commands.Register) {
	So(recordedEvents, ShouldHaveLength, 1)
	registered, ok := recordedEvents[0].(events.Registered)
	So(ok, ShouldBeTrue)
	So(registered.CustomerID().Equals(register.CustomerID()), ShouldBeTrue)
	So(registered.EmailAddress().Equals(register.EmailAddress()), ShouldBeTrue)
	So(registered.ConfirmationHash().Equals(register.ConfirmationHash()), ShouldBeTrue)
	So(registered.PersonName().Equals(register.PersonName()), ShouldBeTrue)
	So(registered.StreamVersion(), ShouldEqual, uint(1))
}
