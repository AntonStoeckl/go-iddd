package valueobjects

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateConfirmationHash(t *testing.T) {
	Convey("Given some input", t, func() {
		input := "foo@bar.com"

		Convey("When GenerateConfirmationHash is invoked", func() {
			confirmationHash := GenerateConfirmationHash(input)

			Convey("Then it should create a ConfirmationHash", func() {
				So(confirmationHash, ShouldImplement, (*ConfirmationHash)(nil))
			})

			Convey("And then it should expose the generated ConfirmationHash", func() {
				So(confirmationHash.String(), ShouldNotBeNil)
			})
		})
	})
}

func TestReconstituteConfirmationHash(t *testing.T) {
	Convey("Given some input", t, func() {
		input := "foo@bar.com"

		Convey("When ReconstituteConfirmationHash is invoked", func() {
			confirmationHash := ReconstituteConfirmationHash(input)

			Convey("Then it should create a ConfirmationHash", func() {
				So(confirmationHash, ShouldImplement, (*ConfirmationHash)(nil))
			})

			Convey("And then it should expose the expected ConfirmationHash", func() {
				So(confirmationHash.String(), ShouldEqual, input)
			})
		})
	})
}

func TestMustMatchOnConfirmationHash(t *testing.T) {
	Convey("Given that two ConfirmationHashs represent equal values", t, func() {
		confirmationHash := ReconstituteConfirmationHash("secret_hash")
		equalConfirmationHash := ReconstituteConfirmationHash("secret_hash")

		Convey("When MustMatch is invoked", func() {
			err := confirmationHash.MustMatch(equalConfirmationHash)

			Convey("Then they should match", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given that two ConfirmationHashs represent different values", t, func() {
		confirmationHash := ReconstituteConfirmationHash("secret_hash")
		differentConfirmationHash := ReconstituteConfirmationHash("different_hash")

		Convey("When MustMatch is invoked", func() {
			err := confirmationHash.MustMatch(differentConfirmationHash)

			Convey("Then they should not match", func() {
				So(err, ShouldBeError, "confirmationHash - is not equal")
			})
		})
	})
}
