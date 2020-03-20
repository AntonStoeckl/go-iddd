package commands_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildRegisterWithInvalidInput(t *testing.T) {
	Convey("When a RegisterCustomer command is built with an invalid emailAddress", t, func() {
		_, err := commands.BuildRegisterCustomer("foo@bar", "John", "Doe")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a RegisterCustomer command is built with an empty givenName", t, func() {
		_, err := commands.BuildRegisterCustomer("john@doe.com", "", "Doe")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a RegisterCustomer command is built with an empty familyName", t, func() {
		_, err := commands.BuildRegisterCustomer("john@doe.com", "John", "")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
