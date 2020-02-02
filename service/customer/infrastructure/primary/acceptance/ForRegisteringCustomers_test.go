package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_ForRegisteringCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)

		commandHandler := application.NewCommandHandler(customerEventStore)

		register, err := commands.NewRegister(
			"john@doe.com",
			"John",
			"Doe",
		)
		So(err, ShouldBeNil)

		containsOnlyRegisteredEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.Registered)

			return ok
		}

		Convey("Given no Customer with the same ID exists", func() {
			customerEventStore.
				On(
					"CreateStreamFrom",
					mock.MatchedBy(containsOnlyRegisteredEvent),
					register.CustomerID(),
				).
				Return(nil).
				Once()

			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Given a Customer with the same ID already exists", func() {
			customerEventStore.
				On(
					"CreateStreamFrom",
					mock.MatchedBy(containsOnlyRegisteredEvent),
					register.CustomerID(),
				).
				Return(lib.ErrDuplicate).
				Once()

			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
				})
			})
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
