package commands_test

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegisterWithInvalidInput(t *testing.T) {
	Convey("When a new Register command is created with an invalid emailAddress", t, func() {
		_, err := commands.NewRegister("foo@bar", "John", "Doe")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a new Register command is created with an empty givenName", t, func() {
		_, err := commands.NewRegister("john@doe.com", "", "Doe")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a new Register command is created with an empty familyName", t, func() {
		_, err := commands.NewRegister("john@doe.com", "John", "")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
