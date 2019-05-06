package valueobjects_test

import (
	"fmt"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestNewConfirmableEmailAddress(t *testing.T) {
	Convey("Given that the supplied emailAddress is valid", t, func() {
		validEmailAddressValue := "foo@bar.com"

		Convey("When NewConfirmableEmailAddress is invoked", func() {
			confirmableEmailAddress, err := valueobjects.NewConfirmableEmailAddress(validEmailAddressValue)

			Convey("Then it should create a ConfirmableEmailAddress", func() {
				So(err, ShouldBeNil)
				So(confirmableEmailAddress, ShouldNotBeNil)
				So(confirmableEmailAddress, ShouldHaveSameTypeAs, (*valueobjects.ConfirmableEmailAddress)(nil))
			})

			// TODO: move all "And ..." into previous block and set FailureMode to FailureHalts (generic?)
			// can be done via init: https://github.com/smartystreets/goconvey/wiki/FAQ
			Convey("And then it should expose the expected value", func() {
				So(confirmableEmailAddress.EmailAddress(), ShouldEqual, validEmailAddressValue)
			})

			Convey("And then it should not be confirmed", func() {
				So(confirmableEmailAddress.IsConfirmed(), ShouldBeFalse)
				// TODO: check if confirmationHash is as expected (can only be done via calling Confirm()
			})
		})
	})

	Convey("Given that the supplied emailAddress is not valid", t, func() {
		invalidEmailAddressValue := "foo@bar.c"

		Convey("When NewConfirmableEmailAddress is invoked", func() {
			confirmableEmailAddress, err := valueobjects.NewConfirmableEmailAddress(invalidEmailAddressValue)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrInvalidInput), ShouldBeTrue)
				So(confirmableEmailAddress, ShouldBeNil)
			})
		})
	})
}

func TestReconstituteConfirmableEmailAddress(t *testing.T) {
	Convey("When ReconstituteConfirmableEmailAddress is invoked", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"
		confirmableEmailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

		Convey("Then it should reconstitute a ConfirmableEmailAddress", func() {
			So(confirmableEmailAddress, ShouldNotBeNil)
			So(confirmableEmailAddress, ShouldHaveSameTypeAs, (*valueobjects.ConfirmableEmailAddress)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(confirmableEmailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})
}

func TestUnmarshalConfirmableEmailAddress(t *testing.T) {
	Convey("Given valid input", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"

		input := map[string]interface{}{
			"emailAddress":     emailAddressValue,
			"confirmationHash": confirmationHashValue,
		}

		Convey("When ReconstituteConfirmableEmailAddress is invoked", func() {
			confirmableEmailAddress, err := valueobjects.UnmarshalConfirmableEmailAddress(input)

			Convey("Then it should unmarshal a ConfirmableEmailAddress", func() {
				So(err, ShouldBeNil)
				So(confirmableEmailAddress, ShouldNotBeNil)
				So(confirmableEmailAddress, ShouldHaveSameTypeAs, (*valueobjects.ConfirmableEmailAddress)(nil))

				Convey("And then it should expose the expected value", func() {
					So(confirmableEmailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
				})
			})
		})
	})

	Convey("Given invalid input", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := 12345 // string expected

		input := map[string]interface{}{
			"emailAddress":     emailAddressValue,
			"confirmationHash": confirmationHashValue,
		}

		Convey("When ReconstituteConfirmableEmailAddress is invoked", func() {
			_, err := valueobjects.UnmarshalConfirmableEmailAddress(input)
			fmt.Println(err)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrUnmarshaling), ShouldBeTrue)
			})
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
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrInvalidInput), ShouldBeTrue)
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

func TestEqualsAnyOnConfirmableEmailAddress(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		confirmationHashValue := "secret_hash"
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress(emailAddressValue, confirmationHashValue)

		Convey("And given another equal EmailAddress", func() {
			equalEmailAddress := valueobjects.ReconstituteEmailAddress(emailAddressValue)

			Convey("When EqualsAny is invoked", func() {
				isEqual := emailAddress.EqualsAny(equalEmailAddress)

				Convey("Then they should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another different EmailAddress", func() {
			differentEmailAddressValue := "foo+different@bar.com"
			differentEmailAddress := valueobjects.ReconstituteEmailAddress(differentEmailAddressValue)

			Convey("When EqualsAny is invoked", func() {
				isEqual := emailAddress.EqualsAny(differentEmailAddress)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
