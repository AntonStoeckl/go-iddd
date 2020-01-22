package test_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ForConfirmingEmailAddresses(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := infrastructure.SetUpDIContainer()
		commandHandler := application.NewCommandHandler(diContainer.GetCustomerRepository(), diContainer.GetPostgresDBConn())

		Convey("Given a registered Customer", func() {
			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)
			So(err, ShouldBeNil)

			Convey("When a Customer's emailAddress is confirmed with a valid confirmationHash", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					register.CustomerID().ID(),
					register.EmailAddress().EmailAddress(),
					register.ConfirmationHash().Hash(),
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("Then it should succeed", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When a Customer's emailAddress is confirmed with an invalid confirmationHash", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					register.CustomerID().ID(),
					register.EmailAddress().EmailAddress(),
					"some_invalid_hash",
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)
				})
			})

			cleanUpArtefactsForPostgresEventStoreSession(diContainer, register.CustomerID())
		})

		Convey("When a Customer's emailAddress is confirmed with an invalid command", func() {
			confirmEmailAddress := commands.ConfirmEmailAddress{}

			err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})
	})
}
