package valueobjects_test

import (
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateConfirmationHash(t *testing.T) {
	Convey("Given some input", t, func() {
		confirmationHashValue := "secret_hash"

		Convey("When GenerateConfirmationHash is invoked", func() {
			confirmationHash := valueobjects.GenerateConfirmationHash(confirmationHashValue)

			Convey("Then it should create a ConfirmationHash", func() {
				So(confirmationHash, ShouldNotBeNil)
				So(confirmationHash, ShouldImplement, (*valueobjects.ConfirmationHash)(nil))
			})

			Convey("And then it should expose the generated ConfirmationHash", func() {
				So(confirmationHash.Hash(), ShouldNotBeNil)
			})
		})
	})
}

func TestReconstituteConfirmationHash(t *testing.T) {
	Convey("When ReconstituteConfirmationHash is invoked", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

		Convey("Then it should reconstitute a ConfirmationHash", func() {
			So(confirmationHash, ShouldNotBeNil)
			So(confirmationHash, ShouldImplement, (*valueobjects.ConfirmationHash)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(confirmationHash.Hash(), ShouldEqual, confirmationHashValue)
		})
	})
}

func TestMustMatchOnConfirmationHash(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

		Convey("And given another equal ConfirmationHash", func() {
			equalConfirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

			Convey("When MustMatch is invoked", func() {
				err := confirmationHash.MustMatch(equalConfirmationHash)

				Convey("Then they must match", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("And given another different ConfirmationHash", func() {
			differentConfirmationHashValue := "different_hash"
			differentConfirmationHash := valueobjects.ReconstituteConfirmationHash(differentConfirmationHashValue)

			Convey("When MustMatch is invoked", func() {
				err := confirmationHash.MustMatch(differentConfirmationHash)

				Convey("Then they must not match", func() {
					So(err, ShouldBeError, "confirmationHash - is not equal")
				})
			})
		})
	})
}
