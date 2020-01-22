package test_test

import (
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ForRegisteringCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := infrastructure.SetUpDIContainer()
		commandHandler := application.NewCommandHandler(
			diContainer.GetCustomerRepository(),
			diContainer.GetPostgresDBConn(),
		)

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

				Convey("And when a Customer is registered with the same ID", func() {
					err = commandHandler.Register(register)

					Convey("It should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
						So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeFalse)
					})
				})
			})

			cleanUpArtefactsForPostgresEventStoreSession(diContainer, register.CustomerID())
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

// TODO: this sucks!
func cleanUpArtefactsForPostgresEventStoreSession(diContainer *cmd.DIContainer, customerID values.CustomerID) {
	streamID := lib.NewStreamID("customer" + "-" + customerID.ID())
	err := diContainer.GetPostgresEventStore().PurgeEventStream(streamID)
	So(err, ShouldBeNil)
}
