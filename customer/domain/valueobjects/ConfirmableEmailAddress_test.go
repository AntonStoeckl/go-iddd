package valueobjects_test

import (
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfirmableEmailAddress(t *testing.T) {
	Convey("Given that the supplied emailAddress is valid", t, func() {
		validEmailAddressValue := "foo@bar.com"

		Convey("When NewConfirmableEmailAddress is invoked", func() {
			confirmableEmailAddress, err := valueobjects.NewConfirmableEmailAddress(validEmailAddressValue)

			Convey("Then it should create a ConfirmableEmailAddress", func() {
				So(err, ShouldBeNil)
				So(confirmableEmailAddress, ShouldNotBeNil)
				So(confirmableEmailAddress, ShouldImplement, (*valueobjects.ConfirmableEmailAddress)(nil))
			})

			Convey("And then it should expose the expected value", func() {
				So(confirmableEmailAddress.EmailAddress(), ShouldEqual, validEmailAddressValue)
			})

			Convey("And then it should not be confirmed", func() {
				So(confirmableEmailAddress.IsConfirmed(), ShouldBeFalse)
			})
		})
	})

	Convey("Given that the supplied emailAddress is not valid", t, func() {
		invalidEmailAddressValue := "foo@bar.c"

		Convey("When NewConfirmableEmailAddress is invoked", func() {
			confirmableEmailAddress, err := valueobjects.NewConfirmableEmailAddress(invalidEmailAddressValue)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "emailAddress - invalid input given")
				So(confirmableEmailAddress, ShouldBeNil)
			})
		})
	})
}

func TestReconstituteConfirmableEmailAddress(t *testing.T) {
	Convey("When ReconstituteConfirmableEmailAddress invoked", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"
		confirmableEmailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

		Convey("Then it should reconstitute a ConfirmableEmailAddress", func() {
			So(confirmableEmailAddress, ShouldNotBeNil)
			So(confirmableEmailAddress, ShouldImplement, (*valueobjects.ConfirmableEmailAddress)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(confirmableEmailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})
}

func TestConfirmOnConfirmableEmailAddress(t *testing.T) {
	Convey("Given an unconfirmed EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"
		differentConfirmationHashValue := "different_secret_hash"
		unconfirmedEmailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

		Convey("When Confirm is invoked with a matching ConfirmationHash", func() {
			suppliedConfirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)
			confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(suppliedConfirmationHash)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
				So(confirmedEmailAddress.Equals(unconfirmedEmailAddress), ShouldBeTrue)
			})

			Convey("And then it should be confirmed", func() {
				So(confirmedEmailAddress.IsConfirmed(), ShouldBeTrue)
			})

			Convey("And then the original should be unchanged", func() {
				So(unconfirmedEmailAddress.IsConfirmed(), ShouldBeFalse)
			})

		})

		Convey("When Confirm is invoked with a different ConfirmationHash", func() {
			given := valueobjects.ReconstituteConfirmationHash(differentConfirmationHashValue)
			confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(given)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "confirmationHash - is not equal")
				So(confirmedEmailAddress, ShouldBeNil)
			})
		})
	})
}

func TestEqualsOnConfirmableEmailAddress(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

		Convey("And given another equal ConfirmableEmailAddress", func() {
			equalEmailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

			Convey("When Equals is invoked", func() {
				isEqual := emailAddress.Equals(equalEmailAddress)

				Convey("Then they should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another different ConfirmableEmailAddress", func() {
			differentEmailAddressValue := "foo+different@bar.com"
			differentEmailAddress := valueobjects.ReconstituteConfirmableEmailAddress(differentEmailAddressValue, confirmationHashValue)

			Convey("When Equals is invoked", func() {
				isEqual := emailAddress.Equals(differentEmailAddress)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
