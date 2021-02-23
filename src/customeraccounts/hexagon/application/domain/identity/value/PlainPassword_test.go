package value_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlainPassword(t *testing.T) {
	Convey("When BuildPassword with valid input", t, func() {
		input := "superSecureP"
		pwd, err := value.BuildPlainPassword(input)

		Convey("Then it should succeed", func() {
			So(err, ShouldBeNil)

			Convey("And the PlainPassword should contain the supplied value", func() {
				So(pwd.String(), ShouldEqual, input)
			})
		})
	})

	Convey("When BuildPassword with leading and trailing", t, func() {
		input := "superSecureP"
		inputWithWhitespace := "\n" + input + " " + "\t"
		pwd, err := value.BuildPlainPassword(inputWithWhitespace)

		Convey("Then it should succeed", func() {
			So(err, ShouldBeNil)

			Convey("And the PlainPassword should contain the value without the whitespace", func() {
				So(pwd.String(), ShouldEqual, input)
			})
		})
	})

	Convey("When RebuildPassword", t, func() {
		input := "superSecurePW123"
		pwd, err := value.RebuildPlainPassword(input)

		Convey("Then it should succeed", func() {
			So(err, ShouldBeNil)

			Convey("And the PlainPassword should contain the supplied value", func() {
				So(pwd.String(), ShouldEqual, input)
			})
		})
	})

	Convey("When BuildPassword with empty input", t, func() {
		input := ""
		_, err := value.BuildPlainPassword(input)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When BuildPassword with too short input", t, func() {
		// 11 chars
		input := "insecure123"
		_, err := value.BuildPlainPassword(input)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When BuildPassword with too long input", t, func() {
		// 251 chars
		input := "hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_hac" +
			"kattack_hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_" +
			"hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_hackattack_hackattac"
		_, err := value.BuildPlainPassword(input)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
