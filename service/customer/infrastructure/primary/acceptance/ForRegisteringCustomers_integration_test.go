package acceptance_test

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Register(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := infrastructure.SetUpDIContainer()
		So(err, ShouldBeNil)
		commandHandler := diContainer.GetCustomerCommandHandler()
		customerEventStore := diContainer.GetCustomerEventStore()

		register, err := commands.NewRegister(
			"john@doe.com",
			"John",
			"Doe",
		)
		So(err, ShouldBeNil)

		Convey("When a Customer is registered", func() {
			err = commandHandler.Register(register)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				eventStream, err := customerEventStore.EventStreamFor(register.CustomerID())
				So(err, ShouldBeNil)
				So(eventStream, ShouldHaveLength, 1)
				So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})

				Convey("And when a Customer is registered with the same ID", func() {
					err = commandHandler.Register(register)

					Convey("It should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
						So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeFalse)
					})
				})
			})

			err := customerEventStore.Delete(register.CustomerID())
			So(err, ShouldBeNil)
		})

		Convey("When a Customer is registered with an invalid command", func() {
			err := commandHandler.Register(commands.Register{})

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})
	})
}
