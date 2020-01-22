package test_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ForChangingEmailAddresses(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := infrastructure.SetUpDIContainer()
		commandHandler := application.NewCommandHandler(
			diContainer.GetCustomerRepository(),
			diContainer.GetPostgresDBConn(),
		)

		newEmailAddress := "john+changed@doe.com"

		Convey("Given a registered Customer", func() {
			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)
			So(err, ShouldBeNil)

			Convey("When a Customer's emailAddress is changed", func() {
				changeEmailAddress, err := commands.NewChangeEmailAddress(
					register.CustomerID().ID(),
					newEmailAddress,
				)
				So(err, ShouldBeNil)

				err = commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("Then it should succeed", func() {
					So(err, ShouldBeNil)

					Convey("When a Customer's emailAddress is changed again to the same emailAddress", func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							newEmailAddress,
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("Then it should succeed", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})

			Convey("When a Customer's emailAddress is changed with an invalid command", func() {
				err := commandHandler.ChangeEmailAddress(commands.ChangeEmailAddress{})

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			cleanUpArtefactsForPostgresEventStoreSession(diContainer, register.CustomerID())
		})

		Convey("Given an unregistered Customer", func() {
			Convey("When a Customer's emailAddress is changed", func() {
				changeEmailAddress, err := commands.NewChangeEmailAddress(
					values.GenerateCustomerID().ID(),
					newEmailAddress,
				)
				So(err, ShouldBeNil)

				err = commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}
